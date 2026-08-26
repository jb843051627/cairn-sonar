package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type FileArchive struct{}

func (FileArchive) Write(ctx context.Context, s model.Survey, e []model.EchoProfile) (model.Archive, error) {
	if err := contextReady(ctx); err != nil {
		return model.Archive{}, err
	}
	h := sha256.New()
	for _, item := range e {
		_, _ = h.Write([]byte(item.ID))
		_, _ = h.Write([]byte(fmt.Sprint(item.DistanceM)))
	}
	digest := hex.EncodeToString(h.Sum(nil))
	return model.Archive{ID: "archive-" + s.ID, SurveyID: s.ID, ObjectKey: "surveys/" + s.ID + "/echoes", Digest: digest, SizeBytes: int64(len(e)), CompletedAt: time.Now().UTC(), Verified: true}, nil
}

func (s *Service) CloseSurvey(ctx context.Context, id string) error {
	survey, err := s.GetSurvey(ctx, id)
	if err != nil {
		return err
	}
	if err := rules.SurveyTransition(survey.Status, model.SurveyClosed); err != nil {
		return err
	}
	survey.Status = model.SurveyClosed
	survey.ClosedAt = utcNow(s.now)
	return s.repo.UpdateSurvey(ctx, survey)
}

func (s *Service) ArchiveSurvey(ctx context.Context, id string) error {
	survey, err := s.GetSurvey(ctx, id)
	if err != nil {
		return err
	}
	if !rules.ArchiveReady(survey) {
		return ErrArchiveBlocked
	}
	echoes, err := s.ListEchoes(ctx, id)
	if err != nil {
		return err
	}
	writer := s.archive
	if writer == nil {
		writer = FileArchive{}
	}
	archive, err := writer.Write(ctx, survey, echoes)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrArchiveFailed, err)
	}
	if err := s.repo.SaveArchive(ctx, archive); err != nil {
		return err
	}
	survey.Status = model.SurveyArchived
	return s.repo.UpdateSurvey(ctx, survey)
}

func (s *Service) GetArchive(ctx context.Context, id string) (model.Archive, error) {
	a, err := s.repo.GetArchive(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Archive{}, nil
	}
	return a, err
}
