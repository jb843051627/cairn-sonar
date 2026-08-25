package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cairn-sonar/internal/model"
)

func (r *Repository) SaveReport(ctx context.Context, report model.SurveyReport) error {
	sections, err := marshal(report.Sections)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO survey_reports(id,survey_id,title,status,summary,sections,generated_at,published_at,version,checksum)
        VALUES(?,?,?,?,?,?,?,?,?,?)
        ON CONFLICT(id) DO UPDATE SET title=excluded.title,status=excluded.status,summary=excluded.summary,
        sections=excluded.sections,generated_at=excluded.generated_at,published_at=excluded.published_at,
        version=excluded.version,checksum=excluded.checksum`,
		report.ID, report.SurveyID, report.Title, report.Status, report.Summary.SurveyID,
		sections, timeText(report.GeneratedAt), timeText(report.PublishedAt), report.Version, report.Checksum)
	if err != nil {
		return fmt.Errorf("save report: %w", err)
	}
	return nil
}

func (r *Repository) GetReport(ctx context.Context, id string) (model.SurveyReport, error) {
	var report model.SurveyReport
	var summary, sections, generated, published string
	err := r.db.QueryRowContext(ctx, `
        SELECT id,survey_id,title,status,summary,sections,generated_at,published_at,version,checksum
        FROM survey_reports WHERE id=?`, id).Scan(
		&report.ID, &report.SurveyID, &report.Title, &report.Status, &summary, &sections,
		&generated, &published, &report.Version, &report.Checksum)
	if errors.Is(err, sql.ErrNoRows) {
		return model.SurveyReport{}, ErrNotFound
	}
	if err != nil {
		return model.SurveyReport{}, fmt.Errorf("get report: %w", err)
	}
	report.GeneratedAt = parseTime(generated)
	report.PublishedAt = parseTime(published)
	if err := unmarshal(sections, &report.Sections); err != nil {
		return model.SurveyReport{}, err
	}
	report.Summary.SurveyID = summary
	return report, nil
}

func (r *Repository) ListReports(ctx context.Context, surveyID, status string, limit, offset int) ([]model.SurveyReport, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	query := `SELECT id,survey_id,title,status,summary,sections,generated_at,published_at,version,checksum FROM survey_reports WHERE survey_id=?`
	args := []any{surveyID}
	if status != "" {
		query += " AND status=?"
		args = append(args, status)
	}
	query += " ORDER BY generated_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.SurveyReport, 0)
	for rows.Next() {
		var report model.SurveyReport
		var summary, sections, generated, published string
		if err := rows.Scan(&report.ID, &report.SurveyID, &report.Title, &report.Status, &summary, &sections, &generated, &published, &report.Version, &report.Checksum); err != nil {
			return nil, err
		}
		report.Summary.SurveyID = summary
		report.GeneratedAt = parseTime(generated)
		report.PublishedAt = parseTime(published)
		if err := unmarshal(sections, &report.Sections); err != nil {
			return nil, err
		}
		out = append(out, report)
	}
	return out, rows.Err()
}

func (r *Repository) DeleteReport(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM survey_reports WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("delete report: %w", err)
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

func (r *Repository) ReportCount(ctx context.Context, surveyID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM survey_reports WHERE survey_id=?", surveyID).Scan(&count)
	return count, err
}
