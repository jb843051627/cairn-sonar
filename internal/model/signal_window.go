package model

import (
	"math"
	"sort"
)

// SampleWindow is the bounded section of a pulse used by the quality rules.
type SampleWindow struct {
	Start       int
	End         int
	Samples     []float64
	SampleRate  float64
	WindowKind  string
	Channel     string
	BaselineDB  float64
	PeakDB      float64
	RMS         float64
	NoiseFloor  float64
	Clipped     bool
	QualityNote string
}

func (w SampleWindow) Clone() SampleWindow {
	w.Samples = w.Samples
	return w
}

func (w SampleWindow) Length() int {
	if w.End <= w.Start {
		return len(w.Samples)
	}
	return w.End - w.Start
}

func (w SampleWindow) DurationSeconds() float64 {
	if w.SampleRate <= 0 {
		return 0
	}
	return float64(w.Length()) / w.SampleRate
}

func (w SampleWindow) Valid() bool {
	return w.Start >= 0 && w.End >= w.Start && len(w.Samples) > 0 && w.SampleRate > 0
}

func (w SampleWindow) Mean() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range w.Samples {
		total += value
	}
	return total / float64(len(w.Samples))
}

func (w SampleWindow) Variance() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	mean := w.Mean()
	total := 0.0
	for _, value := range w.Samples {
		delta := value - mean
		total += delta * delta
	}
	return total / float64(len(w.Samples))
}

func (w SampleWindow) StandardDeviation() float64 {
	return math.Sqrt(w.Variance())
}

func (w SampleWindow) Peak() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	peak := w.Samples[0]
	for _, value := range w.Samples[1:] {
		if value > peak {
			peak = value
		}
	}
	return peak
}

func (w SampleWindow) Trough() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	trough := w.Samples[0]
	for _, value := range w.Samples[1:] {
		if value < trough {
			trough = value
		}
	}
	return trough
}

func (w SampleWindow) CrestFactor() float64 {
	rms := w.RMS
	if rms == 0 {
		rms = math.Sqrt(w.MeanSquare())
	}
	if rms == 0 {
		return 0
	}
	return math.Abs(w.Peak()) / rms
}

func (w SampleWindow) MeanSquare() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range w.Samples {
		total += value * value
	}
	return total / float64(len(w.Samples))
}

func (w SampleWindow) Percentile(percent float64) float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	values := cloneFloats(w.Samples)
	sort.Float64s(values)
	if percent <= 0 {
		return values[0]
	}
	if percent >= 1 {
		return values[len(values)-1]
	}
	position := percent * float64(len(values)-1)
	low := int(position)
	high := low + 1
	if high >= len(values) {
		return values[low]
	}
	fraction := position - float64(low)
	return values[low] + (values[high]-values[low])*fraction
}

func (w SampleWindow) InterquartileRange() float64 {
	return w.Percentile(0.75) - w.Percentile(0.25)
}

func (w SampleWindow) ClippedRatio(limit float64) float64 {
	if len(w.Samples) == 0 || limit <= 0 {
		return 0
	}
	count := 0
	for _, value := range w.Samples {
		if math.Abs(value) >= limit {
			count++
		}
	}
	return float64(count) / float64(len(w.Samples))
}

type SignalFeature struct {
	WindowID      string
	PeakDB        float64
	FloorDB       float64
	ContrastDB    float64
	RMSDB         float64
	CrestFactor   float64
	ZeroCrossings int
	RiseSamples   int
	FallSamples   int
	Energy        float64
	Confidence    float64
	Band          string
	Flags         []string
}

func (f SignalFeature) Clone() SignalFeature {
	f.Flags = append([]string(nil), f.Flags...)
	return f
}

func (f SignalFeature) Reliable() bool {
	return f.ContrastDB >= 6 && f.Confidence >= 0.7 && len(f.Flags) == 0
}

func (f SignalFeature) HasFlag(flag string) bool {
	for _, current := range f.Flags {
		if current == flag {
			return true
		}
	}
	return false
}

func (f *SignalFeature) AddFlag(flag string) {
	if flag == "" || f.HasFlag(flag) {
		return
	}
	f.Flags = append(f.Flags, flag)
}

type SurveyQuality struct {
	SurveyID       string
	PulseCount     int
	UsableCount    int
	ClippedCount   int
	MissingCount   int
	MeanConfidence float64
	Coverage       float64
	Score          float64
	Band           string
	Warnings       []string
}

func (q SurveyQuality) Clone() SurveyQuality {
	q.Warnings = append([]string(nil), q.Warnings...)
	return q
}

func (q SurveyQuality) Complete() bool {
	return q.PulseCount > 0 && q.MissingCount == 0 && q.Coverage >= 0.95 && q.Score >= 0.7
}

func (q *SurveyQuality) AddWarning(warning string) {
	if warning == "" {
		return
	}
	for _, current := range q.Warnings {
		if current == warning {
			return
		}
	}
	q.Warnings = append(q.Warnings, warning)
}

type FeatureBatch struct {
	SurveyID string
	Items    []SignalFeature
	Quality  SurveyQuality
}

func (b FeatureBatch) Clone() FeatureBatch {
	items := make([]SignalFeature, len(b.Items))
	for i := range b.Items {
		items[i] = b.Items[i].Clone()
	}
	b.Items = items
	b.Quality = b.Quality.Clone()
	return b
}

func (b FeatureBatch) CountReliable() int {
	count := 0
	for _, item := range b.Items {
		if item.Reliable() {
			count++
		}
	}
	return count
}

func (b FeatureBatch) MeanContrast() float64 {
	if len(b.Items) == 0 {
		return 0
	}
	total := 0.0
	for _, item := range b.Items {
		total += item.ContrastDB
	}
	return total / float64(len(b.Items))
}
