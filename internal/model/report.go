package model

import "time"

type ReportSection struct {
	Key         string
	Title       string
	Body        string
	Rank        int
	Visible     bool
	FindingIDs  []string
	GeneratedAt time.Time
}

func (s ReportSection) Clone() ReportSection {
	s.FindingIDs = append([]string(nil), s.FindingIDs...)
	return s
}

func (s ReportSection) Valid() bool {
	return s.Key != "" && s.Title != "" && s.Body != "" && s.Rank >= 0
}

type SurveyReport struct {
	ID          string
	SurveyID    string
	Title       string
	Status      string
	Summary     SurveyDigest
	Sections    []ReportSection
	GeneratedAt time.Time
	PublishedAt time.Time
	Version     int
	Checksum    string
}

func (r SurveyReport) Clone() SurveyReport {
	r.Summary = r.Summary.Clone()
	r.Sections = make([]ReportSection, len(r.Sections))
	for i := range r.Sections {
		r.Sections[i] = r.Sections[i].Clone()
	}
	return r
}

func (r SurveyReport) Valid() bool {
	if r.ID == "" || r.SurveyID == "" || r.Title == "" || r.Version < 1 {
		return false
	}
	for _, section := range r.Sections {
		if !section.Valid() {
			return false
		}
	}
	return true
}

func (r SurveyReport) Published() bool {
	return r.Status == "published" && !r.PublishedAt.IsZero()
}

func (r *SurveyReport) AddSection(section ReportSection) {
	if section.Valid() {
		r.Sections = append(r.Sections, section)
	}
}

func (r *SurveyReport) Publish(at time.Time) {
	r.Status = "published"
	r.PublishedAt = at.UTC()
	if r.Version < 1 {
		r.Version = 1
	}
}

func (r *SurveyReport) Unpublish() {
	r.Status = "draft"
	r.PublishedAt = time.Time{}
}

type ReportFilter struct {
	Status       string
	MinSeverity  int
	OnlyCritical bool
	From         time.Time
	To           time.Time
	Limit        int
	Offset       int
}

func (f ReportFilter) Normalize() ReportFilter {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 25
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return f
}

func (f ReportFilter) HasWindow() bool {
	return !f.From.IsZero() || !f.To.IsZero()
}

type ReportIndex struct {
	ReportID  string
	SurveyID  string
	Token     string
	Kind      string
	Value     string
	Weight    float64
	CreatedAt time.Time
	Stale     bool
}

func (i ReportIndex) Valid() bool {
	return i.ReportID != "" && i.SurveyID != "" && i.Token != "" && i.Kind != ""
}

func (i ReportIndex) Clone() ReportIndex { return i }

type ExportBundle struct {
	Survey      Survey
	Report      SurveyReport
	Features    FeatureBatch
	Calibration CalibrationSeries
	GeneratedAt time.Time
	Format      string
	Bytes       int64
}

func (b ExportBundle) Clone() ExportBundle {
	b.Survey = b.Survey.Clone()
	b.Report = b.Report.Clone()
	b.Features = b.Features.Clone()
	b.Calibration = b.Calibration.Clone()
	return b
}

func (b ExportBundle) Ready() bool {
	return b.Survey.ID != "" && b.Report.Valid() && b.GeneratedAt.Unix() > 0 && b.Format != ""
}
