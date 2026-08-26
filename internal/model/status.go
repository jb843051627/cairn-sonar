package model

import "fmt"

type SurveyStatus string

const (
	SurveyPlanned  SurveyStatus = "planned"
	SurveyActive   SurveyStatus = "active"
	SurveyPaused   SurveyStatus = "paused"
	SurveyClosed   SurveyStatus = "closed"
	SurveyArchived SurveyStatus = "archived"
)

func (s SurveyStatus) Valid() bool {
	switch s {
	case SurveyPlanned, SurveyActive, SurveyPaused, SurveyClosed, SurveyArchived:
		return true
	default:
		return false
	}
}

func (s SurveyStatus) String() string { return string(s) }

func CanSurveyTransition(from, to SurveyStatus) bool {
	switch from {
	case SurveyPlanned:
		return to == SurveyActive
	case SurveyActive:
		return to == SurveyPaused || to == SurveyClosed
	case SurveyPaused:
		return to == SurveyActive || to == SurveyClosed
	case SurveyClosed:
		return to == SurveyArchived
	default:
		return false
	}
}

func ValidateSurveyTransition(from, to SurveyStatus) error {
	if !CanSurveyTransition(from, to) {
		return fmt.Errorf("survey transition %s -> %s is not allowed", from, to)
	}
	return nil
}
