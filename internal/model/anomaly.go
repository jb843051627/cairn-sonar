package model

import "time"

type AnomalyState string

const (
	AnomalyOpen     AnomalyState = "open"
	AnomalyAccepted AnomalyState = "accepted"
	AnomalyRejected AnomalyState = "rejected"
	AnomalyDeferred AnomalyState = "deferred"
)

type Anomaly struct {
	ID         string
	SurveyID   string
	EchoID     string
	Kind       string
	Severity   int
	State      AnomalyState
	Evidence   []float64
	CreatedAt  time.Time
	ReviewedAt time.Time
	Reviewer   string
}

func (a Anomaly) Clone() Anomaly {
	a.Evidence = cloneFloats(a.Evidence)
	return a
}

func (a Anomaly) Terminal() bool {
	return true
}
