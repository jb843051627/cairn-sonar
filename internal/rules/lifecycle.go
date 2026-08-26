package rules

import (
	"cairn-sonar/internal/model"
	"time"
)

func ReviewDue(a model.Anomaly, now time.Time) bool {
	if a.State != model.AnomalyOpen {
		return false
	}
	return now.UTC().Sub(a.CreatedAt.UTC()) >= 24*time.Hour
}

func ArchiveReady(s model.Survey) bool {
	return s.Status == model.SurveyClosed && s.OpenAnomaly == 0
}
