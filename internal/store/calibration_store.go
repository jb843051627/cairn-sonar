package store

import (
	"cairn-sonar/internal/model"
	"context"
)

func (r *Repository) SaveCalibration(ctx context.Context, c model.Calibration) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO calibrations(id,instrument_id,survey_id,reference_db,measured_db,offset_db,operator,passed,created_at) VALUES(?,?,?,?,?,?,?,?,?)`, c.ID, c.InstrumentID, c.SurveyID, c.ReferenceDB, c.MeasuredDB, c.OffsetDB, c.Operator, boolInt(c.Passed), timeText(c.CreatedAt))
	return err
}
func (r *Repository) CountCalibrations(ctx context.Context, instrumentID string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM calibrations WHERE instrument_id=?`, instrumentID).Scan(&n)
	return n, err
}
