package rules

import "cairn-sonar/internal/model"

type Thresholds struct {
	PeakDB, NoiseDB, Confidence float64
	SeverityDB                  float64
}

func DefaultThresholds() Thresholds {
	return Thresholds{PeakDB: -18, NoiseDB: -58, Confidence: .72, SeverityDB: 12}
}

func (t Thresholds) Score(e model.EchoProfile) (float64, int) {
	contrast := e.PeakDB - e.NoiseDB
	score := e.Confidence * contrast / 40
	severity := 0
	if contrast >= t.SeverityDB {
		severity = 3
	} else if contrast >= 8 {
		severity = 2
	} else if contrast >= 5 {
		severity = 1
	}
	if score > 1 {
		score = 1
	}
	if score < 0 {
		score = 0
	}
	return score, severity
}
