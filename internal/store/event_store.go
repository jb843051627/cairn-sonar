package store

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) AppendEvent(ctx context.Context, surveyID, kind, payload string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO events(survey_id,kind,payload,created_at) VALUES(?,?,?,?)`, surveyID, kind, payload, timeText(at))
	if err != nil {
		return fmt.Errorf("append event: %w", err)
	}
	return nil
}
func (r *Repository) EventCount(ctx context.Context, surveyID string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM events WHERE survey_id=?`, surveyID).Scan(&n)
	return n, err
}
