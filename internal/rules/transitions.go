package rules

import (
	"cairn-sonar/internal/model"
	"fmt"
)

func ReviewTransition(from, to model.AnomalyState) error {
	if from != model.AnomalyOpen {
		return fmt.Errorf("anomaly %s cannot be reviewed", from)
	}
	switch to {
	case model.AnomalyAccepted, model.AnomalyRejected, model.AnomalyDeferred:
		return nil
	default:
		return fmt.Errorf("unsupported review state %s", to)
	}
}

func SurveyTransition(from, to model.SurveyStatus) error {
	return model.ValidateSurveyTransition(from, to)
}
