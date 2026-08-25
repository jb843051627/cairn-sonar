package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"context"
	"sort"
	"time"
)

func (s *Service) DueAnomalies(ctx context.Context, surveyID string) ([]model.Anomaly, error) {
	items, err := s.ListAnomalies(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	now := s.now.Now()
	sort.SliceStable(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	out := make([]model.Anomaly, 0)
	for _, a := range items {
		if rules.ReviewDue(a, now) {
			out = append(out, a)
		}
	}
	return out, nil
}
func (s *Service) ReviewWindow(a model.Anomaly) time.Duration {
	if a.CreatedAt.IsZero() {
		return 0
	}
	return s.now.Now().Sub(a.CreatedAt)
}
