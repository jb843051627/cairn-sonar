package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cairn-sonar/internal/model"
)

func (r *Repository) SaveEnvironment(ctx context.Context, surveyID string, snapshot model.EnvironmentSnapshot) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO environments(survey_id,temperature_c,humidity_pct,pressure_hpa,recorded_at) VALUES(?,?,?,?,?)`, surveyID, snapshot.TemperatureC, snapshot.HumidityPct, snapshot.PressureHPA, timeText(snapshot.RecordedAt))
	if err != nil {
		return fmt.Errorf("save environment: %w", err)
	}
	return nil
}

func (r *Repository) LatestEnvironment(ctx context.Context, surveyID string) (model.EnvironmentSnapshot, error) {
	var snapshot model.EnvironmentSnapshot
	var recorded string
	err := r.db.QueryRowContext(ctx, `SELECT temperature_c,humidity_pct,pressure_hpa,recorded_at FROM environments WHERE survey_id=? ORDER BY recorded_at DESC LIMIT 1`, surveyID).Scan(&snapshot.TemperatureC, &snapshot.HumidityPct, &snapshot.PressureHPA, &recorded)
	if errors.Is(err, sql.ErrNoRows) {
		return model.EnvironmentSnapshot{}, ErrNotFound
	}
	if err != nil {
		return model.EnvironmentSnapshot{}, fmt.Errorf("latest environment: %w", err)
	}
	snapshot.RecordedAt = parseTime(recorded)
	return snapshot, nil
}

func (r *Repository) EnvironmentCount(ctx context.Context, surveyID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM environments WHERE survey_id=?", surveyID).Scan(&count)
	return count, err
}

func (r *Repository) DeleteEnvironment(ctx context.Context, surveyID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM environments WHERE survey_id=?", surveyID)
	return err
}
