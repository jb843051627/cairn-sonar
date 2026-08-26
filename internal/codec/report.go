package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cairn-sonar/internal/model"
)

var ErrInvalidReport = errors.New("invalid report document")

type ReportDocument struct {
	Title    string               `json:"title"`
	SurveyID string               `json:"survey_id"`
	Status   string               `json:"status"`
	Summary  string               `json:"summary"`
	Sections []ReportDocumentPart `json:"sections"`
}

type ReportDocumentPart struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Visible bool   `json:"visible"`
}

func EncodeReport(report model.SurveyReport) ([]byte, error) {
	document := ReportDocument{Title: report.Title, SurveyID: report.SurveyID, Status: report.Status, Summary: strings.Join(report.Summary.Highlights, "；"), Sections: make([]ReportDocumentPart, 0, len(report.Sections))}
	for _, section := range report.Sections {
		document.Sections = append(document.Sections, ReportDocumentPart{Key: section.Key, Title: section.Title, Body: section.Body, Visible: section.Visible})
	}
	return json.MarshalIndent(document, "", "  ")
}

func DecodeReport(data []byte) (ReportDocument, error) {
	var document ReportDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return ReportDocument{}, err
	}
	if document.Title == "" || document.SurveyID == "" || len(document.Sections) == 0 {
		return ReportDocument{}, ErrInvalidReport
	}
	return document, nil
}

func RenderReportText(document ReportDocument) string {
	var builder strings.Builder
	builder.WriteString(document.Title)
	builder.WriteString("\n勘测批次：")
	builder.WriteString(document.SurveyID)
	builder.WriteString("\n状态：")
	builder.WriteString(document.Status)
	builder.WriteString("\n")
	if document.Summary != "" {
		builder.WriteString("摘要：")
		builder.WriteString(document.Summary)
		builder.WriteString("\n")
	}
	for _, section := range document.Sections {
		if !section.Visible {
			continue
		}
		builder.WriteString("\n## ")
		builder.WriteString(section.Title)
		builder.WriteString("\n")
		builder.WriteString(section.Body)
		builder.WriteString("\n")
	}
	return builder.String()
}

func ValidateReportDocument(document ReportDocument) error {
	if document.Title == "" {
		return fmt.Errorf("title is required")
	}
	if document.SurveyID == "" {
		return fmt.Errorf("survey id is required")
	}
	if len(document.Sections) == 0 {
		return fmt.Errorf("at least one report section is required")
	}
	for index, section := range document.Sections {
		if section.Key == "" || section.Title == "" || section.Body == "" {
			return fmt.Errorf("section %d is incomplete", index)
		}
	}
	return nil
}
