package model

import "time"

type Review struct {
	ID        string
	AnomalyID string
	Decision  AnomalyState
	Comment   string
	Reviewer  string
	CreatedAt time.Time
}

func (r Review) Valid() bool {
	return r.ID != "" && r.AnomalyID != "" && r.Reviewer != "" &&
		(r.Decision == AnomalyAccepted || r.Decision == AnomalyRejected || r.Decision == AnomalyDeferred)
}
