package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) SaveArchive(ctx context.Context, a model.Archive) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO archives(id,survey_id,object_key,digest,size_bytes,completed_at,verified) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET digest=excluded.digest,size_bytes=excluded.size_bytes,completed_at=excluded.completed_at,verified=excluded.verified`, a.ID, a.SurveyID, a.ObjectKey, a.Digest, a.SizeBytes, timeText(a.CompletedAt), boolInt(a.Verified))
	return err
}
func (r *Repository) GetArchive(ctx context.Context, surveyID string) (model.Archive, error) {
	var a model.Archive
	var completed string
	err := r.db.QueryRowContext(ctx, `SELECT id,survey_id,object_key,digest,size_bytes,completed_at,verified FROM archives WHERE survey_id=? ORDER BY completed_at DESC LIMIT 1`, surveyID).Scan(&a.ID, &a.SurveyID, &a.ObjectKey, &a.Digest, &a.SizeBytes, &completed, &a.Verified)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Archive{}, ErrNotFound
	}
	if err != nil {
		return model.Archive{}, fmt.Errorf("get archive: %w", err)
	}
	a.CompletedAt = parseTime(completed)
	return a, nil
}
