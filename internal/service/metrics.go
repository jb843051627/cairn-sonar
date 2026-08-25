package service

import "context"

type Metrics struct{ Pulses, Echoes, Anomalies, Events int }

func (s *Service) Metrics(ctx context.Context, surveyID string) (Metrics, error) {
	survey, err := s.GetSurvey(ctx, surveyID)
	if err != nil {
		return Metrics{}, err
	}
	events, err := s.repo.EventCount(ctx, surveyID)
	if err != nil {
		return Metrics{}, err
	}
	anomalies, err := s.repo.ListAnomalies(ctx, surveyID)
	if err != nil {
		return Metrics{}, err
	}
	return Metrics{Pulses: survey.PulseCount, Echoes: survey.EchoCount, Anomalies: len(anomalies), Events: events}, nil
}
