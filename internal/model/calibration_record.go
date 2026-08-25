package model

import "time"

type CalibrationRecord struct {
	ID             string
	InstrumentID   string
	SurveyID       string
	Operator       string
	ReferenceDB    float64
	MeasuredDB     float64
	NoiseDB        float64
	FrequencyHz    float64
	TemperatureC   float64
	HumidityPct    float64
	OffsetDB       float64
	DriftPpm       float64
	ToleranceDB    float64
	Passed         bool
	RetestRequired bool
	CreatedAt      time.Time
	Notes          []string
}

func (r CalibrationRecord) Clone() CalibrationRecord {
	r.Notes = append([]string(nil), r.Notes...)
	return r
}

func (r CalibrationRecord) ErrorDB() float64 { return r.MeasuredDB - r.ReferenceDB }

func (r CalibrationRecord) WithinTolerance() bool {
	tolerance := r.ToleranceDB
	if tolerance <= 0 {
		tolerance = 2
	}
	return absFloat(r.ErrorDB()) <= tolerance
}

func (r CalibrationRecord) Stable() bool {
	return r.WithinTolerance() && absFloat(r.DriftPpm) < 30 && r.NoiseDB < -35
}

func (r *CalibrationRecord) Evaluate() {
	r.OffsetDB = r.ReferenceDB - r.MeasuredDB
	r.Passed = r.WithinTolerance() && r.Stable()
	r.RetestRequired = !r.Passed
}

func (r *CalibrationRecord) AddNote(note string) {
	if note != "" {
		r.Notes = append(r.Notes, note)
	}
}

type CalibrationSeries struct {
	InstrumentID string
	Records      []CalibrationRecord
	MeanErrorDB  float64
	MaxErrorDB   float64
	Stable       bool
}

func (s CalibrationSeries) Clone() CalibrationSeries {
	s.Records = make([]CalibrationRecord, len(s.Records))
	for i := range s.Records {
		s.Records[i] = s.Records[i].Clone()
	}
	return s
}

func (s CalibrationSeries) PassedCount() int {
	count := 0
	for _, record := range s.Records {
		if record.Passed {
			count++
		}
	}
	return count
}

type EnvironmentSnapshot struct {
	TemperatureC float64
	HumidityPct  float64
	PressureHPA  float64
	RecordedAt   time.Time
}

func (e EnvironmentSnapshot) Valid() bool {
	return e.TemperatureC > -50 && e.TemperatureC < 80 && e.HumidityPct >= 0 && e.HumidityPct <= 100 && e.PressureHPA > 0
}

func (e EnvironmentSnapshot) AcousticFactor() float64 {
	temperature := 331.3 + 0.606*e.TemperatureC
	humidity := 1 - (e.HumidityPct-50)/1000
	if humidity < 0.8 {
		humidity = 0.8
	}
	return temperature * humidity
}

func (e EnvironmentSnapshot) Clone() EnvironmentSnapshot { return e }

type InstrumentHealth struct {
	InstrumentID string
	Samples      int
	Failures     int
	LastSuccess  time.Time
	LastFailure  time.Time
	DriftPpm     float64
	Enabled      bool
	Reasons      []string
}

func (h InstrumentHealth) Clone() InstrumentHealth {
	h.Reasons = append([]string(nil), h.Reasons...)
	return h
}

func (h InstrumentHealth) FailureRate() float64 {
	if h.Samples <= 0 {
		return 0
	}
	return float64(h.Failures) / float64(h.Samples)
}

func (h InstrumentHealth) Healthy() bool {
	return h.Enabled && h.FailureRate() < 0.05 && absFloat(h.DriftPpm) < 30
}

func (h *InstrumentHealth) AddReason(reason string) {
	if reason != "" {
		h.Reasons = append(h.Reasons, reason)
	}
}

func absFloat(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
