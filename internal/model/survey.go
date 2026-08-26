package model

import "time"

type Survey struct {
	ID          string
	ChamberID   string
	Lead        string
	Status      SurveyStatus
	StartedAt   time.Time
	ClosedAt    time.Time
	PulseCount  int
	EchoCount   int
	OpenAnomaly int
	Notes       []string
}

func (s Survey) Clone() Survey {
	s.Notes = cloneStrings(s.Notes)
	return s
}

func (s Survey) Valid() bool {
	return s.ID != "" && s.ChamberID != "" && s.Lead != "" && s.Status.Valid()
}

func (s Survey) CanAcceptData() bool {
	return s.Status == SurveyActive || s.Status == SurveyPaused
}

func (s *Survey) AddNote(note string) {
	if note != "" {
		s.Notes = append(s.Notes, note)
	}
}
