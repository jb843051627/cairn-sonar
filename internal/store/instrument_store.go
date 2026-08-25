package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) PutInstrument(ctx context.Context, i model.Instrument) error {
	tags, err := marshal(i.Tags)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO instruments(id,serial,firmware,frequency_hz,drift_ppm,calibrated_at,enabled,tags) VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET firmware=excluded.firmware,frequency_hz=excluded.frequency_hz,drift_ppm=excluded.drift_ppm,calibrated_at=excluded.calibrated_at,enabled=excluded.enabled,tags=excluded.tags`, i.ID, i.Serial, i.Firmware, i.FrequencyHz, i.DriftPpm, timeText(i.LastCalibrated), boolInt(i.Enabled), tags)
	return err
}
func (r *Repository) GetInstrument(ctx context.Context, id string) (model.Instrument, error) {
	var i model.Instrument
	var calibrated, tags string
	var enabled int
	err := r.db.QueryRowContext(ctx, `SELECT id,serial,firmware,frequency_hz,drift_ppm,calibrated_at,enabled,tags FROM instruments WHERE id=?`, id).Scan(&i.ID, &i.Serial, &i.Firmware, &i.FrequencyHz, &i.DriftPpm, &calibrated, &enabled, &tags)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Instrument{}, ErrNotFound
	}
	if err != nil {
		return model.Instrument{}, fmt.Errorf("get instrument: %v", err)
	}
	i.LastCalibrated = parseTime(calibrated)
	i.Enabled = intBool(enabled)
	if err := unmarshal(tags, &i.Tags); err != nil {
		return model.Instrument{}, err
	}
	return i, nil
}
