package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) InsertAnomaly(ctx context.Context, a model.Anomaly) error {
	e, err := marshal(a.Evidence)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO anomalies(id,survey_id,echo_id,kind,severity,state,evidence,created_at,reviewed_at,reviewer) VALUES(?,?,?,?,?,?,?,?,?,?)`, a.ID, a.SurveyID, a.EchoID, a.Kind, a.Severity, a.State, e, timeText(a.CreatedAt), timeText(a.ReviewedAt), a.Reviewer)
	return err
}
func (r *Repository) GetAnomaly(ctx context.Context, id string) (model.Anomaly, error) {
	var a model.Anomaly
	var evidence, created, reviewed string
	err := r.db.QueryRowContext(ctx, `SELECT id,survey_id,echo_id,kind,severity,state,evidence,created_at,reviewed_at,reviewer FROM anomalies WHERE id=?`, id).Scan(&a.ID, &a.SurveyID, &a.EchoID, &a.Kind, &a.Severity, &a.State, &evidence, &created, &reviewed, &a.Reviewer)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Anomaly{}, fmt.Errorf("get anomaly: %v", ErrNotFound)
	}
	if err != nil {
		return model.Anomaly{}, fmt.Errorf("get anomaly: %v", err)
	}
	a.CreatedAt = parseTime(created)
	a.ReviewedAt = parseTime(reviewed)
	if err := unmarshal(evidence, &a.Evidence); err != nil {
		return model.Anomaly{}, err
	}
	return a, nil
}
func (r *Repository) UpdateAnomaly(ctx context.Context, a model.Anomaly) error {
	e, err := marshal(a.Evidence)
	if err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx, `UPDATE anomalies SET state=?,evidence=?,reviewed_at=?,reviewer=? WHERE id=?`, a.State, e, timeText(a.ReviewedAt), a.Reviewer, a.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *Repository) ListAnomalies(ctx context.Context, surveyID string) ([]model.Anomaly, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,survey_id,echo_id,kind,severity,state,evidence,created_at,reviewed_at,reviewer FROM anomalies WHERE survey_id=? ORDER BY severity DESC,created_at`, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Anomaly, 0)
	for rows.Next() {
		var a model.Anomaly
		var e, c, rv string
		if err := rows.Scan(&a.ID, &a.SurveyID, &a.EchoID, &a.Kind, &a.Severity, &a.State, &e, &c, &rv, &a.Reviewer); err != nil {
			return nil, err
		}
		a.CreatedAt = parseTime(c)
		a.ReviewedAt = parseTime(rv)
		if err := unmarshal(e, &a.Evidence); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
