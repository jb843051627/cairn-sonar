package service

import (
	"context"
	"fmt"
)

func (s *Service) ReconcileSurvey(ctx context.Context, id string) error {
	survey, err := s.GetSurvey(ctx, id)
	if err != nil {
		return err
	}
	pulses, err := s.ListPulses(ctx, id)
	if err != nil {
		return err
	}
	echoes, err := s.ListEchoes(ctx, id)
	if err != nil {
		return err
	}
	if survey.PulseCount != len(pulses) || survey.EchoCount != len(echoes) {
		return fmt.Errorf("survey counters drift: pulses=%d/%d echoes=%d/%d", survey.PulseCount, len(pulses), survey.EchoCount, len(echoes))
	}
	return nil
}
