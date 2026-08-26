package model

import "time"

type EchoProfile struct {
	ID         string
	SurveyID   string
	PulseID    string
	PeakDB     float64
	NoiseDB    float64
	DistanceM  float64
	Confidence float64
	Bands      []float64
	CreatedAt  time.Time
}

func (e EchoProfile) Clone() EchoProfile {
	e.Bands = cloneFloats(e.Bands)
	return e
}

func (e EchoProfile) Reliable() bool {
	return e.Confidence >= 0.75 && e.PeakDB-e.NoiseDB >= 6
}
