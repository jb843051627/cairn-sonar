package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"cairn-sonar/internal/model"
	"cairn-sonar/internal/store"
)

func (s *Service) BuildReport(ctx context.Context, surveyID string) (model.SurveyReport, error) {
	survey, err := s.GetSurvey(ctx, surveyID)
	if err != nil {
		return model.SurveyReport{}, err
	}
	quality, _, err := s.Signals().Quality(ctx, surveyID)
	if err != nil {
		return model.SurveyReport{}, err
	}
	findings, err := s.ListFindings(ctx, surveyID, 1)
	if err != nil {
		return model.SurveyReport{}, err
	}
	pulses, err := s.ListPulses(ctx, surveyID)
	if err != nil {
		return model.SurveyReport{}, err
	}
	echoes, err := s.ListEchoes(ctx, surveyID)
	if err != nil {
		return model.SurveyReport{}, err
	}
	digest := model.SurveyDigest{SurveyID: surveyID, Status: survey.Status, PulseCount: len(pulses), EchoCount: len(echoes), FindingCount: len(findings), QualityScore: quality.Score, GeneratedAt: utcNow(s.now)}
	for _, finding := range findings {
		if finding.Severity >= 3 {
			digest.CriticalCount++
			digest.Highlight(finding.Category + " 需要重点复核")
		}
	}
	if quality.Complete() {
		digest.Highlight("声学采样质量完整")
	} else {
		digest.Highlight("采样质量需要人工复核")
	}
	report := model.SurveyReport{ID: "report-" + surveyID, SurveyID: surveyID, Title: "石窟岩壁声学巡检报告", Status: "draft", Summary: digest, GeneratedAt: utcNow(s.now), Version: 1}
	report.AddSection(model.ReportSection{Key: "summary", Title: "巡检摘要", Body: fmt.Sprintf("共采集 %d 个脉冲，生成 %d 个回波，质量得分 %.2f。", len(pulses), len(echoes), quality.Score), Rank: 1, Visible: true, GeneratedAt: utcNow(s.now)})
	report.AddSection(model.ReportSection{Key: "findings", Title: "异常复核", Body: findingSummary(findings), Rank: 2, Visible: true, FindingIDs: findingIDs(findings), GeneratedAt: utcNow(s.now)})
	report.AddSection(model.ReportSection{Key: "quality", Title: "数据质量", Body: strings.Join(digest.Highlights, "；"), Rank: 3, Visible: true, GeneratedAt: utcNow(s.now)})
	report.Checksum = reportChecksum(report)
	if err := s.repo.SaveReport(ctx, report); err != nil {
		return model.SurveyReport{}, err
	}
	return report.Clone(), nil
}

func (s *Service) GetReport(ctx context.Context, id string) (model.SurveyReport, error) {
	report, err := s.repo.GetReport(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return model.SurveyReport{}, ErrNoEcho
	}
	if err != nil {
		return model.SurveyReport{}, err
	}
	return report.Clone(), nil
}

func (s *Service) ListReports(ctx context.Context, surveyID, status string, limit, offset int) ([]model.SurveyReport, error) {
	reports, err := s.repo.ListReports(ctx, surveyID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]model.SurveyReport, len(reports))
	for i := range reports {
		out[i] = reports[i].Clone()
	}
	return out, nil
}

func (s *Service) PublishReport(ctx context.Context, id string) error {
	report, err := s.GetReport(ctx, id)
	if err != nil {
		return err
	}
	if !report.Valid() {
		return fmt.Errorf("report is incomplete")
	}
	report.Publish(utcNow(s.now))
	if err := s.repo.SaveReport(ctx, report); err != nil {
		return err
	}
	return s.repo.AppendEvent(ctx, report.SurveyID, "report.published", report.ID, utcNow(s.now))
}

func (s *Service) UnpublishReport(ctx context.Context, id string) error {
	report, err := s.GetReport(ctx, id)
	if err != nil {
		return err
	}
	report.Unpublish()
	return s.repo.SaveReport(ctx, report)
}

func (s *Service) IndexReport(ctx context.Context, report model.SurveyReport) (int, error) {
	count := 0
	for _, section := range report.Sections {
		for _, token := range strings.Fields(section.Body) {
			if token == "" {
				continue
			}
			index := model.ReportIndex{ReportID: report.ID, SurveyID: report.SurveyID, Token: token, Kind: section.Key, Value: section.Title, Weight: float64(section.Rank), CreatedAt: utcNow(s.now)}
			if err := s.repo.SaveReportIndex(ctx, index); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

func (s *Service) SearchReports(ctx context.Context, surveyID, token string, limit int) ([]model.ReportIndex, error) {
	return s.repo.SearchReportIndex(ctx, surveyID, token, limit)
}

func (s *Service) DeleteReport(ctx context.Context, id string) error {
	return s.repo.DeleteReport(ctx, id)
}

func findingSummary(findings []model.Finding) string {
	if len(findings) == 0 {
		return "未发现需要复核的异常。"
	}
	parts := make([]string, 0, len(findings))
	for _, finding := range findings {
		parts = append(parts, fmt.Sprintf("%s(级别%d)", finding.Category, finding.Severity))
	}
	return strings.Join(parts, "、")
}

func findingIDs(findings []model.Finding) []string {
	ids := make([]string, 0, len(findings))
	for _, finding := range findings {
		ids = append(ids, finding.ID)
	}
	return ids
}

func reportChecksum(report model.SurveyReport) string {
	h := sha256.New()
	_, _ = h.Write([]byte(report.ID))
	_, _ = h.Write([]byte(report.Title))
	for _, section := range report.Sections {
		_, _ = h.Write([]byte(section.Key))
		_, _ = h.Write([]byte(section.Body))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func reportAge(report model.SurveyReport, now time.Time) time.Duration {
	if report.GeneratedAt.IsZero() {
		return 0
	}
	return now.Sub(report.GeneratedAt)
}
