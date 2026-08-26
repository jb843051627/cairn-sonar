package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) UpsertEcho(ctx context.Context, e model.EchoProfile) error {
	bands, err := marshal(e.Bands)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO echoes(id,survey_id,pulse_id,peak_db,noise_db,distance_m,confidence,bands,created_at) VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET peak_db=excluded.peak_db,noise_db=excluded.noise_db,distance_m=excluded.distance_m,confidence=excluded.confidence,bands=excluded.bands`, e.ID, e.SurveyID, e.PulseID, e.PeakDB, e.NoiseDB, e.DistanceM, e.Confidence, bands, timeText(e.CreatedAt))
	return err
}

func (r *Repository) GetEcho(ctx context.Context, id string) (model.EchoProfile, error) {
	var e model.EchoProfile
	var bands, created string
	err := r.db.QueryRowContext(ctx, `SELECT id,survey_id,pulse_id,peak_db,noise_db,distance_m,confidence,bands,created_at FROM echoes WHERE id=?`, id).Scan(&e.ID, &e.SurveyID, &e.PulseID, &e.PeakDB, &e.NoiseDB, &e.DistanceM, &e.Confidence, &bands, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return model.EchoProfile{}, ErrNotFound
	}
	if err != nil {
		return model.EchoProfile{}, fmt.Errorf("get echo: %w", err)
	}
	e.CreatedAt = parseTime(created)
	if err := unmarshal(bands, &e.Bands); err != nil {
		return model.EchoProfile{}, err
	}
	return e, nil
}

func (r *Repository) ListEchoes(ctx context.Context, surveyID string) ([]model.EchoProfile, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,survey_id,pulse_id,peak_db,noise_db,distance_m,confidence,bands,created_at FROM echoes WHERE survey_id=? ORDER BY distance_m`, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.EchoProfile, 0)
	for rows.Next() {
		var e model.EchoProfile
		var bands, created string
		if err := rows.Scan(&e.ID, &e.SurveyID, &e.PulseID, &e.PeakDB, &e.NoiseDB, &e.DistanceM, &e.Confidence, &bands, &created); err != nil {
			return nil, err
		}
		e.CreatedAt = parseTime(created)
		if err := unmarshal(bands, &e.Bands); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
