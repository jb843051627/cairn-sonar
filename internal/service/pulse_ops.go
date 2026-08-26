package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/store"
	"context"
	"errors"
	"fmt"
)

func (s *Service) RecordPulse(ctx context.Context, p model.Pulse) error {
	if err := contextReady(ctx); err != nil {
		return err
	}
	if !p.Valid() {
		return ErrInvalidPulse
	}
	survey, err := s.GetSurvey(ctx, p.SurveyID)
	if err != nil {
		return err
	}
	if !survey.CanAcceptData() {
		return fmt.Errorf("survey %s is not accepting pulses", p.SurveyID)
	}
	if _, err := s.repo.GetInstrument(ctx, p.InstrumentID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return ErrInstrumentNotFound
		}
		return err
	}
	if err := s.repo.InsertPulse(ctx, p.Clone()); err != nil {
		return err
	}
	survey.PulseCount++
	if err := s.repo.UpdateSurvey(ctx, survey); err != nil {
		return err
	}
	if err := contextReady(ctx); err != nil {
		return err
	}
	s.queue.Submit(p.SurveyID)
	return s.repo.AppendEvent(ctx, p.SurveyID, "pulse.recorded", p.ID, utcNow(s.now))
}

func (s *Service) IngestBatch(ctx context.Context, surveyID string, pulses []model.Pulse) (int, error) {
	if err := contextReady(ctx); err != nil {
		return 0, err
	}
	accepted := 0
	for _, p := range pulses {
		p.SurveyID = surveyID
		if err := contextReady(ctx); err != nil {
			return accepted, err
		}
		if err := s.RecordPulse(ctx, p); err != nil {
			return accepted, err
		}
		accepted++
	}
	return accepted, nil
}

func (s *Service) ListPulses(ctx context.Context, surveyID string) ([]model.Pulse, error) {
	return s.repo.ListPulses(ctx, surveyID)
}
