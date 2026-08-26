package service

import (
	"context"
	"encoding/json"
)

type Export struct {
	Survey    any `json:"survey"`
	Pulses    any `json:"pulses"`
	Echoes    any `json:"echoes"`
	Anomalies any `json:"anomalies"`
}

func (s *Service) ExportSurvey(ctx context.Context, id string) ([]byte, error) {
	survey, err := s.GetSurvey(ctx, id)
	if err != nil {
		return nil, err
	}
	pulses, err := s.ListPulses(ctx, id)
	if err != nil {
		return nil, err
	}
	echoes, err := s.ListEchoes(ctx, id)
	if err != nil {
		return nil, err
	}
	anomalies, err := s.ListAnomalies(ctx, id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Export{Survey: survey, Pulses: pulses, Echoes: echoes, Anomalies: anomalies})
}
