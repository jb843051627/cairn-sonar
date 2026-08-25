package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) InsertSurvey(ctx context.Context, s model.Survey) error {
	notes, err := marshal(s.Notes)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO surveys(id,chamber_id,lead,status,started_at,closed_at,pulse_count,echo_count,open_anomaly,notes) VALUES(?,?,?,?,?,?,?,?,?,?)`, s.ID, s.ChamberID, s.Lead, s.Status, timeText(s.StartedAt), timeText(s.ClosedAt), s.PulseCount, s.EchoCount, s.OpenAnomaly, notes)
	if err != nil {
		return fmt.Errorf("insert survey: %w", err)
	}
	return nil
}

func (r *Repository) GetSurvey(ctx context.Context, id string) (model.Survey, error) {
	var s model.Survey
	var started, closed, notes string
	err := r.db.QueryRowContext(ctx, `SELECT id,chamber_id,lead,status,started_at,closed_at,pulse_count,echo_count,open_anomaly,notes FROM surveys WHERE id=?`, id).Scan(&s.ID, &s.ChamberID, &s.Lead, &s.Status, &started, &closed, &s.PulseCount, &s.EchoCount, &s.OpenAnomaly, &notes)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Survey{}, ErrNotFound
	}
	if err != nil {
		return model.Survey{}, fmt.Errorf("get survey: %v", err)
	}
	s.StartedAt = parseTime(started)
	s.ClosedAt = parseTime(closed)
	if err := unmarshal(notes, &s.Notes); err != nil {
		return model.Survey{}, err
	}
	return s, nil
}

func (r *Repository) UpdateSurvey(ctx context.Context, s model.Survey) error {
	notes, err := marshal(s.Notes)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `UPDATE surveys SET chamber_id=?,lead=?,status=?,started_at=?,closed_at=?,pulse_count=?,echo_count=?,open_anomaly=?,notes=? WHERE id=?`, s.ChamberID, s.Lead, s.Status, timeText(s.StartedAt), timeText(s.ClosedAt), s.PulseCount, s.EchoCount, s.OpenAnomaly, notes, s.ID)
	if err != nil {
		return fmt.Errorf("update survey: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ListSurveys(ctx context.Context, status model.SurveyStatus) ([]model.Survey, error) {
	query, args := `SELECT id,chamber_id,lead,status,started_at,closed_at,pulse_count,echo_count,open_anomaly,notes FROM surveys ORDER BY started_at`, []any{}
	if status != "" {
		query += ` WHERE status=?`
		args = append(args, status)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Survey, 0)
	for rows.Next() {
		var s model.Survey
		var a, b, n string
		if err := rows.Scan(&s.ID, &s.ChamberID, &s.Lead, &s.Status, &a, &b, &s.PulseCount, &s.EchoCount, &s.OpenAnomaly, &n); err != nil {
			return nil, err
		}
		s.StartedAt = parseTime(a)
		s.ClosedAt = parseTime(b)
		if err := unmarshal(n, &s.Notes); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
