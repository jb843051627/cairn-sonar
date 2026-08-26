package service

import (
	"context"
	"fmt"
	"sort"

	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
)

type SignalService struct {
	parent *Service
	config rules.AnalysisConfig
}

func (s *Service) Signals() SignalService {
	return SignalService{parent: s, config: rules.DefaultAnalysisConfig()}
}

func (s SignalService) Extract(ctx context.Context, pulse model.Pulse, from, to int) (model.SignalFeature, error) {
	if err := contextReady(ctx); err != nil {
		return model.SignalFeature{}, err
	}
	window, err := rules.DetectWindow(pulse, from, to, s.config)
	if err != nil {
		return model.SignalFeature{}, fmt.Errorf("detect signal window: %w", err)
	}
	return rules.FeatureFromWindow(pulse.ID, window, s.config), nil
}

func (s SignalService) ExtractAll(ctx context.Context, pulses []model.Pulse) ([]model.SignalFeature, error) {
	features := make([]model.SignalFeature, 0, len(pulses))
	for _, pulse := range pulses {
		if err := contextReady(ctx); err != nil {
			return features, err
		}
		feature, err := s.Extract(ctx, pulse, 0, len(pulse.Samples))
		if err != nil {
			return features, err
		}
		features = append(features, feature)
	}
	return features, nil
}

func (s SignalService) Quality(ctx context.Context, surveyID string) (model.SurveyQuality, []rules.RuleResult, error) {
	pulses, err := s.parent.ListPulses(ctx, surveyID)
	if err != nil {
		return model.SurveyQuality{}, nil, err
	}
	features, err := s.ExtractAll(ctx, pulses)
	if err != nil {
		return model.SurveyQuality{}, nil, err
	}
	batch := model.FeatureBatch{SurveyID: surveyID, Items: features}
	quality, results := rules.EvaluateSurveyQuality(batch, rules.DefaultQualityRules())
	return quality, results, nil
}

func (s SignalService) CriticalFeatures(ctx context.Context, surveyID string) ([]model.SignalFeature, error) {
	quality, results, err := s.Quality(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	_ = quality
	failed := rules.FailedRules(results)
	ids := make(map[string]bool)
	for _, result := range failed {
		for _, evidence := range result.Evidence {
			ids[evidence] = true
		}
	}
	pulses, err := s.parent.ListPulses(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	features, err := s.ExtractAll(ctx, pulses)
	if err != nil {
		return nil, err
	}
	out := make([]model.SignalFeature, 0)
	for _, feature := range features {
		if ids[feature.WindowID] || feature.ContrastDB >= 14 {
			out = append(out, feature)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ContrastDB > out[j].ContrastDB })
	return out, nil
}

func (s SignalService) Classify(ctx context.Context, pulse model.Pulse) (string, int, string, error) {
	feature, err := s.Extract(ctx, pulse, 0, len(pulse.Samples))
	if err != nil {
		return "", 0, "", err
	}
	kind, severity, tag := rules.ClassifyFeature(feature, rules.DefaultSeverityRules())
	return kind, severity, tag, nil
}

func (s SignalService) Windows(pulse model.Pulse, width int) []model.SampleWindow {
	if width <= 0 {
		width = len(pulse.Samples)
	}
	out := make([]model.SampleWindow, 0)
	for from := 0; from < len(pulse.Samples); from += width {
		to := from + width
		if to > len(pulse.Samples) {
			to = len(pulse.Samples)
		}
		window, err := rules.DetectWindow(pulse, from, to, s.config)
		if err == nil {
			out = append(out, window)
		}
	}
	return out
}
