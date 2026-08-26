package rules

import (
	"cairn-sonar/internal/model"
	"errors"
	"math"
)

var ErrNoSamples = errors.New("no acoustic samples")

func EchoFromSamples(p model.Pulse) (model.EchoProfile, error) {
	if len(p.Samples) == 0 {
		return model.EchoProfile{}, ErrNoSamples
	}
	peak, sum := p.Samples[0], 0.0
	for _, sample := range p.Samples {
		if sample > peak {
			peak = sample
		}
		sum += sample
	}
	mean := sum / float64(len(p.Samples))
	noise := mean - 3
	variance := 0.0
	for _, sample := range p.Samples {
		d := sample - mean
		variance += d * d
	}
	variance /= float64(len(p.Samples))
	confidence := 1 / (1 + math.Sqrt(variance)/10)
	return model.EchoProfile{ID: "echo-" + p.ID, SurveyID: p.SurveyID, PulseID: p.ID, PeakDB: peak, NoiseDB: noise, DistanceM: float64(p.DurationMS) * 0.17, Confidence: confidence, Bands: []float64{peak, mean, variance}}, nil
}
