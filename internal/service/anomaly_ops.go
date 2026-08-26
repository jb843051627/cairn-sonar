package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"context"
	"fmt"
)

func (s *Service) RaiseAnomaly(ctx context.Context, a model.Anomaly) error {
	if err := contextReady(ctx); err != nil {
		return err
	}
	if a.ID == "" || a.SurveyID == "" || a.EchoID == "" {
		return fmt.Errorf("invalid anomaly")
	}
	a.State = model.AnomalyOpen
	if a.CreatedAt.IsZero() {
		a.CreatedAt = utcNow(s.now)
	}
	if err := s.repo.InsertAnomaly(ctx, a.Clone()); err != nil {
		return err
	}
	survey, err := s.GetSurvey(ctx, a.SurveyID)
	if err != nil {
		return err
	}
	survey.OpenAnomaly++
	if err := s.repo.UpdateSurvey(ctx, survey); err != nil {
		return err
	}
	s.cacheMu.Lock()
	s.anomalyCache[a.SurveyID] = append(s.anomalyCache[a.SurveyID], a.Clone())
	s.cacheMu.Unlock()
	return s.repo.AppendEvent(ctx, a.SurveyID, "anomaly.raised", a.ID, utcNow(s.now))
}

func (s *Service) ReviewAnomaly(ctx context.Context, id string, decision model.AnomalyState, reviewer, comment string) error {
	a, err := s.repo.GetAnomaly(ctx, id)
	if err != nil {
		return s.mapStoreError(err, ErrAnomalyNotFound)
	}
	if err := rules.ReviewTransition(a.State, decision); err != nil {
		return err
	}
	a.State = decision
	a.Reviewer = reviewer
	a.ReviewedAt = utcNow(s.now)
	if err := s.repo.UpdateAnomaly(ctx, a); err != nil {
		return err
	}
	s.cacheMu.Lock()
	if cached, ok := s.anomalyCache[a.SurveyID]; ok {
		for i := range cached {
			if cached[i].ID == a.ID {
				cached[i] = a.Clone()
			}
		}
	}
	s.cacheMu.Unlock()
	survey, err := s.GetSurvey(ctx, a.SurveyID)
	if err != nil {
		return err
	}
	if survey.OpenAnomaly > 0 {
		survey.OpenAnomaly--
	}
	if err := s.repo.UpdateSurvey(ctx, survey); err != nil {
		return err
	}
	return s.repo.AppendEvent(ctx, a.SurveyID, "anomaly.reviewed", comment, utcNow(s.now))
}

func (s *Service) ListAnomalies(ctx context.Context, surveyID string) ([]model.Anomaly, error) {
	s.cacheMu.RLock()
	cached, ok := s.anomalyCache[surveyID]
	s.cacheMu.RUnlock()
	if ok {
		return cloneAnomalies(cached), nil
	}
	items, err := s.repo.ListAnomalies(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	s.cacheMu.Lock()
	s.anomalyCache[surveyID] = cloneAnomalies(items)
	s.cacheMu.Unlock()
	return cloneAnomalies(items), nil
}
func cloneAnomalies(in []model.Anomaly) []model.Anomaly {
	out := make([]model.Anomaly, len(in))
	for i := range in {
		out[i] = in[i].Clone()
	}
	return out
}
