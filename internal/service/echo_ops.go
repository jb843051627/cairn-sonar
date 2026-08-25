package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"cairn-sonar/internal/store"
	"context"
	"errors"
	"fmt"
	"sort"
)

func (s *Service) AnalyzeEcho(ctx context.Context, p model.Pulse) (model.EchoProfile, error) {

	echo, err := rules.EchoFromSamples(p)
	if err != nil {
		return model.EchoProfile{}, fmt.Errorf("decode echo: %w", err)
	}
	echo.CreatedAt = utcNow(s.now)
	if err := s.repo.UpsertEcho(ctx, echo); err != nil {
		return model.EchoProfile{}, err
	}
	survey, err := s.GetSurvey(ctx, p.SurveyID)
	if err == nil {
		survey.PulseCount++
		_ = s.repo.UpdateSurvey(ctx, survey)
	}
	s.cacheMu.Lock()
	s.echoCache[p.SurveyID] = append(s.echoCache[p.SurveyID], echo.Clone())
	s.cacheMu.Unlock()
	return echo, nil
}

func (s *Service) ListEchoes(ctx context.Context, surveyID string) ([]model.EchoProfile, error) {
	s.cacheMu.RLock()
	cached, ok := s.echoCache[surveyID]
	s.cacheMu.RUnlock()
	if ok {
		return cloneEchoes(cached), nil
	}
	items, err := s.repo.ListEchoes(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	s.cacheMu.Lock()
	s.echoCache[surveyID] = cloneEchoes(items)
	s.cacheMu.Unlock()
	return cloneEchoes(items), nil
}

func (s *Service) SortedEchoes(ctx context.Context, surveyID string) ([]model.EchoProfile, error) {
	items, err := s.ListEchoes(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Confidence > items[j].Confidence })
	return items, nil
}

func (s *Service) EchoByID(ctx context.Context, id string) (model.EchoProfile, error) {
	e, err := s.repo.GetEcho(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return model.EchoProfile{}, ErrNoEcho
	}
	return e, err
}
func cloneEchoes(in []model.EchoProfile) []model.EchoProfile {
	out := make([]model.EchoProfile, len(in))
	for i := range in {
		out[i] = in[i].Clone()
	}
	return out
}
