package model

import "time"

type Pulse struct {
	ID           string
	SurveyID     string
	InstrumentID string
	Sequence     int
	EmittedAt    time.Time
	DurationMS   int
	GainDB       float64
	Samples      []float64
	Tags         []string
}

func (p Pulse) Clone() Pulse {
	p.Samples = cloneFloats(p.Samples)
	p.Tags = cloneStrings(p.Tags)
	return p
}

func (p Pulse) Valid() bool {
	return p.ID != "" && p.SurveyID != "" && p.InstrumentID != "" && p.Sequence >= 0 && len(p.Samples) > 0
}
