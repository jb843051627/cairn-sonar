package service

import (
	"cairn-sonar/internal/model"
	"context"
)

func (s *Service) ResumeSurvey(ctx context.Context, id string) error {
	return s.transitionSurvey(ctx, id, model.SurveyActive)
}
func (s *Service) CloseAndArchive(ctx context.Context, id string) error {
	if err := s.CloseSurvey(ctx, id); err != nil {
		return err
	}
	return s.ArchiveSurvey(ctx, id)
}
