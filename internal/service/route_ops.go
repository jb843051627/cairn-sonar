package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"context"
	"fmt"
)

func (s *Service) PlanRoute(ctx context.Context, surveyID string) (model.Route, error) {
	if err := contextReady(ctx); err != nil {
		return model.Route{}, err
	}
	survey, err := s.GetSurvey(ctx, surveyID)
	if err != nil {
		return model.Route{}, err
	}
	chambers, err := s.repo.ListChambers(ctx)
	if err != nil {
		return model.Route{}, err
	}
	if len(chambers) == 0 {
		return model.Route{}, fmt.Errorf("no chambers for survey")
	}
	route := rules.BuildRoute(chambers)
	route.ID = "route-" + survey.ID
	route.SurveyID = survey.ID
	route.CreatedAt = utcNow(s.now)
	if err := s.repo.SaveRoute(ctx, route); err != nil {
		return model.Route{}, err
	}
	return route, nil
}
func (s *Service) GetRoute(ctx context.Context, id string) (model.Route, error) {
	return s.repo.GetRoute(ctx, id)
}
