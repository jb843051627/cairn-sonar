package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"cairn-sonar/internal/store"
)

func (s *Service) CreateFinding(ctx context.Context, finding model.Finding) error {
	if err := contextReady(ctx); err != nil {
		return err
	}
	if !finding.Valid() {
		return fmt.Errorf("invalid finding")
	}
	if finding.State == "" {
		finding.State = model.AnomalyOpen
	}
	if finding.CreatedAt.IsZero() {
		finding.CreatedAt = utcNow(s.now)
	}
	finding.UpdatedAt = finding.CreatedAt
	if err := s.repo.SaveFinding(ctx, finding.Clone()); err != nil {
		return err
	}
	return s.repo.AppendEvent(ctx, finding.SurveyID, "finding.created", finding.ID, utcNow(s.now))
}

func (s *Service) GetFinding(ctx context.Context, id string) (model.Finding, error) {
	finding, err := s.repo.GetFinding(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return model.Finding{}, ErrAnomalyNotFound
	}
	if err != nil {
		return model.Finding{}, err
	}
	return finding.Clone(), nil
}

func (s *Service) ListFindings(ctx context.Context, surveyID string, minSeverity int) ([]model.Finding, error) {
	items, err := s.repo.ListFindings(ctx, surveyID, minSeverity)
	if err != nil {
		return nil, err
	}
	out := make([]model.Finding, len(items))
	for i := range items {
		out[i] = items[i].Clone()
	}
	return out, nil
}

func (s *Service) AutoCreateFinding(ctx context.Context, pulse model.Pulse) (model.Finding, error) {
	kind, severity, tag, err := s.Signals().Classify(ctx, pulse)
	if err != nil {
		return model.Finding{}, err
	}
	if severity == 0 {
		return model.Finding{}, nil
	}
	finding := model.Finding{ID: "finding-" + pulse.ID, SurveyID: pulse.SurveyID, EchoID: "echo-" + pulse.ID, Category: kind, Severity: severity, Confidence: 0.8, Description: "自动识别的" + tag + "声学异常", Evidence: append([]float64(nil), pulse.Samples...)}
	if err := s.CreateFinding(ctx, finding); err != nil {
		return model.Finding{}, err
	}
	return finding, nil
}

func (s *Service) ResolveFinding(ctx context.Context, id string, decision model.AnomalyState, reviewer, reason string) error {
	finding, err := s.GetFinding(ctx, id)
	if err != nil {
		return err
	}
	if err := rules.ReviewTransition(finding.State, decision); err != nil {
		return err
	}
	trail := model.ReviewTrail{ID: id + "-" + string(decision), FindingID: id, From: finding.State, To: decision, Reviewer: reviewer, Reason: reason, CreatedAt: utcNow(s.now), Confidence: finding.Confidence}
	if !trail.Valid() {
		return fmt.Errorf("invalid review trail")
	}
	finding.State = decision
	finding.UpdatedAt = utcNow(s.now)
	if err := s.repo.SaveFinding(ctx, finding); err != nil {
		return err
	}
	return s.repo.SaveReviewTrail(ctx, trail)
}

func (s *Service) DueFindings(ctx context.Context, surveyID string, now time.Time) ([]model.Finding, error) {
	items, err := s.ListFindings(ctx, surveyID, 1)
	if err != nil {
		return nil, err
	}
	out := make([]model.Finding, 0)
	for _, item := range items {
		if item.State == model.AnomalyOpen && now.Sub(item.CreatedAt) >= 24*time.Hour {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Service) OpenFindingCount(ctx context.Context, surveyID string) (int, error) {
	return s.repo.CountOpenFindings(ctx, surveyID)
}
