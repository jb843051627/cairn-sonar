package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cairn-sonar/internal/model"
)

func (r *Repository) SaveReportIndex(ctx context.Context, index model.ReportIndex) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO report_indices(report_id,survey_id,token,kind,value,weight,created_at,stale) VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(report_id,token) DO UPDATE SET value=excluded.value,weight=excluded.weight,stale=excluded.stale`, index.ReportID, index.SurveyID, index.Token, index.Kind, index.Value, index.Weight, timeText(index.CreatedAt), boolInt(index.Stale))
	return err
}

func (r *Repository) SearchReportIndex(ctx context.Context, surveyID, token string, limit int) ([]model.ReportIndex, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	rows, err := r.db.QueryContext(ctx, `SELECT report_id,survey_id,token,kind,value,weight,created_at,stale FROM report_indices WHERE survey_id=? AND token LIKE ? ORDER BY weight DESC LIMIT ?`, surveyID, "%"+token+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.ReportIndex, 0)
	for rows.Next() {
		var index model.ReportIndex
		var created string
		var stale int
		if err := rows.Scan(&index.ReportID, &index.SurveyID, &index.Token, &index.Kind, &index.Value, &index.Weight, &created, &stale); err != nil {
			return nil, err
		}
		index.CreatedAt = parseTime(created)
		index.Stale = intBool(stale)
		out = append(out, index)
	}
	return out, rows.Err()
}

func (r *Repository) MarkIndexStale(ctx context.Context, reportID string) error {
	result, err := r.db.ExecContext(ctx, "UPDATE report_indices SET stale=1 WHERE report_id=?", reportID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) CountIndexEntries(ctx context.Context, reportID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM report_indices WHERE report_id=?", reportID).Scan(&count)
	return count, err
}

func (r *Repository) DeleteStaleIndices(ctx context.Context, before string) (int64, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM report_indices WHERE stale=1 AND created_at<?", before)
	if err != nil {
		return 0, fmt.Errorf("delete stale indices: %w", err)
	}
	return result.RowsAffected()
}

func (r *Repository) ReportExists(ctx context.Context, reportID string) (bool, error) {
	var value int
	err := r.db.QueryRowContext(ctx, "SELECT 1 FROM survey_reports WHERE id=? LIMIT 1", reportID).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return value == 1, nil
}
