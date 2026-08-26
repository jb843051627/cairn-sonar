package model

import "time"

type Calibration struct {
	ID           string
	InstrumentID string
	SurveyID     string
	ReferenceDB  float64
	MeasuredDB   float64
	OffsetDB     float64
	Operator     string
	Passed       bool
	CreatedAt    time.Time
}

func (c Calibration) Delta() float64 { return c.MeasuredDB - c.ReferenceDB }

func (c Calibration) Valid() bool {
	return c.ID != "" && c.InstrumentID != "" && c.Operator != "" && c.ReferenceDB >= 0
}
