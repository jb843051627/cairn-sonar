package rules

import (
	"errors"
	"math"
	"sort"

	"cairn-sonar/internal/model"
)

var ErrInvalidWindow = errors.New("invalid signal window")

type AnalysisConfig struct {
	ClipLimit       float64
	MinContrastDB   float64
	MinConfidence   float64
	BaselineSamples int
	Smoothing       int
}

func DefaultAnalysisConfig() AnalysisConfig {
	return AnalysisConfig{ClipLimit: 95, MinContrastDB: 6, MinConfidence: 0.7, BaselineSamples: 8, Smoothing: 3}
}

func NormalizeSamples(samples []float64) []float64 {
	if len(samples) == 0 {
		return nil
	}
	peak := 0.0
	for _, value := range samples {
		if abs(value) > peak {
			peak = abs(value)
		}
	}
	if peak == 0 {
		return append([]float64(nil), samples...)
	}
	out := make([]float64, len(samples))
	for i, value := range samples {
		out[i] = value / peak * 100
	}
	return out
}

func SmoothSamples(samples []float64, width int) []float64 {
	if len(samples) == 0 {
		return nil
	}
	if width < 2 {
		return samples
	}
	if width > len(samples) {
		width = len(samples)
	}
	out := make([]float64, len(samples))
	half := width / 2
	for i := range samples {
		from := i - half
		to := i + half + 1
		if from < 0 {
			from = 0
		}
		if to > len(samples) {
			to = len(samples)
		}
		total := 0.0
		for _, value := range samples[from:to] {
			total += value
		}
		out[i] = total / float64(to-from)
	}
	return out
}

func DetectWindow(pulse model.Pulse, from, to int, cfg AnalysisConfig) (model.SampleWindow, error) {
	if from < 0 {
		from = 0
	}
	if to <= 0 || to > len(pulse.Samples) {
		to = len(pulse.Samples)
	}
	if from >= to || to > len(pulse.Samples) {
		return model.SampleWindow{}, ErrInvalidWindow
	}
	samples := append([]float64(nil), pulse.Samples[from:to]...)
	if cfg.Smoothing > 1 {
		samples = SmoothSamples(samples, cfg.Smoothing)
	}
	w := model.SampleWindow{Start: from, End: to, Samples: samples, SampleRate: 1000, WindowKind: "echo", Channel: pulse.InstrumentID}
	w.PeakDB = w.Peak()
	w.RMS = math.Sqrt(w.MeanSquare())
	w.NoiseFloor = w.Percentile(0.25)
	w.BaselineDB = w.Mean()
	w.Clipped = w.ClippedRatio(cfg.ClipLimit) > 0.01
	if w.Clipped {
		w.QualityNote = "clipped samples"
	}
	return w, nil
}

func FeatureFromWindow(id string, w model.SampleWindow, cfg AnalysisConfig) model.SignalFeature {
	feature := model.SignalFeature{WindowID: id, PeakDB: w.Peak(), FloorDB: w.NoiseFloor, RMSDB: w.RMS, CrestFactor: w.CrestFactor(), Energy: w.MeanSquare()}
	feature.ContrastDB = feature.PeakDB - feature.FloorDB
	feature.Confidence = confidence(feature, cfg)
	feature.ZeroCrossings = zeroCrossings(w.Samples)
	feature.RiseSamples, feature.FallSamples = edgeLengths(w.Samples)
	feature.Band = bandForPeak(feature.PeakDB)
	if w.Clipped {
		feature.AddFlag("clipped")
	}
	if feature.ContrastDB < cfg.MinContrastDB {
		feature.AddFlag("low-contrast")
	}
	if feature.Confidence < cfg.MinConfidence {
		feature.AddFlag("low-confidence")
	}
	return feature
}

func confidence(feature model.SignalFeature, cfg AnalysisConfig) float64 {
	contrast := feature.ContrastDB / 35
	if contrast < 0 {
		contrast = 0
	}
	if contrast > 1 {
		contrast = 1
	}
	stability := 1 / (1 + feature.CrestFactor/10)
	confidence := contrast*0.7 + stability*0.3
	if confidence > 1 {
		return 1
	}
	if confidence < 0 {
		return 0
	}
	return confidence
}

func zeroCrossings(samples []float64) int {
	if len(samples) < 2 {
		return 0
	}
	count := 0
	previous := samples[0]
	for _, current := range samples[1:] {
		if (previous < 0 && current >= 0) || (previous >= 0 && current < 0) {
			count++
		}
		previous = current
	}
	return count
}

func edgeLengths(samples []float64) (int, int) {
	if len(samples) < 2 {
		return 0, 0
	}
	peak := 0
	for i := range samples {
		if samples[i] > samples[peak] {
			peak = i
		}
	}
	fall := len(samples) - peak - 1
	return peak, fall
}

func bandForPeak(peak float64) string {
	switch {
	case peak >= 80:
		return "high"
	case peak >= 45:
		return "mid"
	default:
		return "low"
	}
}

func RankFeatures(items []model.SignalFeature) []model.SignalFeature {
	out := make([]model.SignalFeature, len(items))
	for i := range items {
		out[i] = items[i].Clone()
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Confidence == out[j].Confidence {
			return out[i].ContrastDB > out[j].ContrastDB
		}
		return out[i].Confidence > out[j].Confidence
	})
	return out
}

func BuildFeatureBatch(surveyID string, windows []model.SampleWindow, cfg AnalysisConfig) model.FeatureBatch {
	items := make([]model.SignalFeature, 0, len(windows))
	for i, window := range windows {
		items = append(items, FeatureFromWindow(surveyID+"-window-"+itoa(i+1), window, cfg))
	}
	quality := model.SurveyQuality{SurveyID: surveyID, PulseCount: len(windows), UsableCount: len(items)}
	for _, item := range items {
		quality.MeanConfidence += item.Confidence
		if item.HasFlag("clipped") {
			quality.ClippedCount++
		}
		if item.Reliable() {
			quality.UsableCount++
		}
	}
	if len(items) > 0 {
		quality.MeanConfidence /= float64(len(items))
		quality.Coverage = float64(quality.UsableCount) / float64(len(items)*2)
	}
	quality.Score = quality.MeanConfidence*0.7 + quality.Coverage*0.3
	quality.Band = qualityBand(quality.Score)
	if quality.ClippedCount > 0 {
		quality.AddWarning("one or more windows contain clipped samples")
	}
	return model.FeatureBatch{SurveyID: surveyID, Items: items, Quality: quality}
}

func qualityBand(score float64) string {
	switch {
	case score >= 0.85:
		return "excellent"
	case score >= 0.7:
		return "good"
	case score >= 0.5:
		return "review"
	default:
		return "poor"
	}
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 8)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
