package model

import "time"

type InspectionLedger struct {
	ID           string
	SurveyID     string
	Operator     string
	StartedAt    time.Time
	FinishedAt   time.Time
	Expected     int
	Observed     int
	Accepted     int
	Rejected     int
	OpenFindings int
	Status       string
	Notes        []string
}

func (l InspectionLedger) Clone() InspectionLedger {
	l.Notes = append([]string(nil), l.Notes...)
	return l
}

func (l InspectionLedger) Progress() float64 {
	if l.Expected <= 0 {
		return 0
	}
	value := float64(l.Observed) / float64(l.Expected)
	if value > 1 {
		return 1
	}
	if value < 0 {
		return 0
	}
	return value
}

func (l InspectionLedger) Balanced() bool {
	return l.Observed == l.Accepted+l.Rejected+l.OpenFindings
}

func (l InspectionLedger) Finished() bool {
	return l.Status == "finished" && !l.FinishedAt.IsZero() && l.OpenFindings == 0
}

func (l *InspectionLedger) Observe(accepted bool) {
	l.Observed++
	if accepted {
		l.Accepted++
	} else {
		l.OpenFindings++
	}
}

func (l *InspectionLedger) ResolveFinding(accepted bool) {
	if l.OpenFindings > 0 {
		l.OpenFindings--
	}
	if accepted {
		l.Accepted++
	} else {
		l.Rejected++
	}
}

func (l *InspectionLedger) AddNote(note string) {
	if note != "" {
		l.Notes = append(l.Notes, note)
	}
}

type Finding struct {
	ID          string
	SurveyID    string
	EchoID      string
	Category    string
	Severity    int
	Confidence  float64
	Description string
	State       AnomalyState
	Evidence    []float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tags        []string
}

func (f Finding) Clone() Finding {
	f.Evidence = cloneFloats(f.Evidence)
	f.Tags = cloneStrings(f.Tags)
	return f
}

func (f Finding) Valid() bool {
	return f.ID != "" && f.SurveyID != "" && f.EchoID != "" && f.Category != "" && f.Severity > 0
}

func (f Finding) Urgent() bool {
	return f.Severity >= 3 && f.Confidence >= 0.8 && f.State == AnomalyOpen
}

func (f *Finding) Tag(tag string) {
	if tag == "" {
		return
	}
	for _, current := range f.Tags {
		if current == tag {
			return
		}
	}
	f.Tags = append(f.Tags, tag)
}

type EvidenceLink struct {
	ID         string
	FindingID  string
	PulseID    string
	SampleFrom int
	SampleTo   int
	Digest     string
	Comment    string
	CreatedAt  time.Time
}

func (e EvidenceLink) Valid() bool {
	return e.ID != "" && e.FindingID != "" && e.PulseID != "" && e.SampleFrom >= 0 && e.SampleTo >= e.SampleFrom
}

func (e EvidenceLink) Width() int {
	return e.SampleTo - e.SampleFrom
}

type ReviewTrail struct {
	ID         string
	FindingID  string
	From       AnomalyState
	To         AnomalyState
	Reviewer   string
	Reason     string
	CreatedAt  time.Time
	Automatic  bool
	Confidence float64
}

func (r ReviewTrail) Valid() bool {
	return r.ID != "" && r.FindingID != "" && r.Reviewer != "" && r.From != r.To && r.Reason != ""
}

type InspectionNote struct {
	ID        string
	SurveyID  string
	Author    string
	Body      string
	Kind      string
	CreatedAt time.Time
	Pinned    bool
}

func (n InspectionNote) Valid() bool {
	return n.ID != "" && n.SurveyID != "" && n.Author != "" && n.Body != ""
}

func (n InspectionNote) Clone() InspectionNote { return n }

type SurveyDigest struct {
	SurveyID      string
	Status        SurveyStatus
	PulseCount    int
	EchoCount     int
	FindingCount  int
	CriticalCount int
	QualityScore  float64
	GeneratedAt   time.Time
	Highlights    []string
}

func (d SurveyDigest) Clone() SurveyDigest {
	d.Highlights = append([]string(nil), d.Highlights...)
	return d
}

func (d SurveyDigest) HasCriticalFinding() bool { return d.CriticalCount > 0 }

func (d *SurveyDigest) Highlight(value string) {
	if value != "" {
		d.Highlights = append(d.Highlights, value)
	}
}
