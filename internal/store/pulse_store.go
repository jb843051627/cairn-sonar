package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) InsertPulse(ctx context.Context, p model.Pulse) error {
	samples, err := marshal(p.Samples)
	if err != nil {
		return err
	}
	tags, err := marshal(p.Tags)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO pulses(id,survey_id,instrument_id,sequence,emitted_at,duration_ms,gain_db,samples,tags) VALUES(?,?,?,?,?,?,?,?,?)`, p.ID, p.SurveyID, p.InstrumentID, p.Sequence, timeText(p.EmittedAt), p.DurationMS, p.GainDB, samples, tags)
	if err != nil {
		return fmt.Errorf("insert pulse: %w", err)
	}
	return nil
}

func (r *Repository) GetPulse(ctx context.Context, id string) (model.Pulse, error) {
	var p model.Pulse
	var emitted, samples, tags string
	err := r.db.QueryRowContext(ctx, `SELECT id,survey_id,instrument_id,sequence,emitted_at,duration_ms,gain_db,samples,tags FROM pulses WHERE id=?`, id).Scan(&p.ID, &p.SurveyID, &p.InstrumentID, &p.Sequence, &emitted, &p.DurationMS, &p.GainDB, &samples, &tags)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Pulse{}, ErrNotFound
	}
	if err != nil {
		return model.Pulse{}, fmt.Errorf("get pulse: %w", err)
	}
	p.EmittedAt = parseTime(emitted)
	if err := unmarshal(samples, &p.Samples); err != nil {
		return model.Pulse{}, err
	}
	if err := unmarshal(tags, &p.Tags); err != nil {
		return model.Pulse{}, err
	}
	return p, nil
}

func (r *Repository) ListPulses(ctx context.Context, surveyID string) ([]model.Pulse, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,survey_id,instrument_id,sequence,emitted_at,duration_ms,gain_db,samples,tags FROM pulses WHERE survey_id=? ORDER BY sequence`, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Pulse, 0)
	for rows.Next() {
		var p model.Pulse
		var emitted, samples, tags string
		if err := rows.Scan(&p.ID, &p.SurveyID, &p.InstrumentID, &p.Sequence, &emitted, &p.DurationMS, &p.GainDB, &samples, &tags); err != nil {
			return nil, err
		}
		p.EmittedAt = parseTime(emitted)
		if err := unmarshal(samples, &p.Samples); err != nil {
			return nil, err
		}
		if err := unmarshal(tags, &p.Tags); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
