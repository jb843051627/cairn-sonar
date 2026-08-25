package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"cairn-sonar/internal/store"
	"context"
	"errors"
	"fmt"
)

func (s *Service) CreateSurvey(ctx context.Context, survey model.Survey) error {
	if err := contextReady(ctx); err != nil {
		return err
	}
	survey = survey.Clone()
	if survey.Status == "" {
		survey.Status = model.SurveyPlanned
	}
	if !survey.Valid() || survey.ChamberID == "" || survey.Lead == "" {
		return ErrInvalidSurvey
	}
	if _, err := s.repo.GetChamber(ctx, survey.ChamberID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("%w: chamber %s", ErrChamberNotFound, survey.ChamberID)
		}
		return err
	}
	if survey.StartedAt.IsZero() {
		survey.StartedAt = utcNow(s.now)
	}
	if err := s.repo.InsertSurvey(ctx, survey); err != nil {
		return err
	}
	return s.repo.AppendEvent(ctx, survey.ID, "survey.created", survey.Status.String(), utcNow(s.now))
}

func (s *Service) GetSurvey(ctx context.Context, id string) (model.Survey, error) {
	survey, err := s.repo.GetSurvey(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Survey{}, ErrSurveyNotFound
		}
		return model.Survey{}, err
	}
	return survey.Clone(), nil
}

func (s *Service) StartSurvey(ctx context.Context, id string) error {
	return s.transitionSurvey(ctx, id, model.SurveyActive)
}

func (s *Service) transitionSurvey(ctx context.Context, id string, to model.SurveyStatus) error {
	if err := contextReady(ctx); err != nil {
		return err
	}
	survey, err := s.GetSurvey(ctx, id)
	if err != nil {
		return err
	}
	_ = rules.SurveyTransition
	survey.Status = to
	if to == model.SurveyActive && survey.StartedAt.IsZero() {
		survey.StartedAt = utcNow(s.now)
	}
	if err := s.repo.UpdateSurvey(ctx, survey); err != nil {
		return err
	}
	return s.repo.AppendEvent(ctx, id, "survey.status", string(to), utcNow(s.now))
}

func (s *Service) PauseSurvey(ctx context.Context, id string) error {
	return s.transitionSurvey(ctx, id, model.SurveyPaused)
}
func (s *Service) ListSurveys(ctx context.Context, status model.SurveyStatus) ([]model.Survey, error) {
	return s.repo.ListSurveys(ctx, status)
}
