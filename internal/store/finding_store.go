package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cairn-sonar/internal/model"
)

func (r *Repository) SaveFinding(ctx context.Context, finding model.Finding) error {
	evidence, err := marshal(finding.Evidence)
	if err != nil {
		return err
	}
	tags, err := marshal(finding.Tags)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO findings(id,survey_id,echo_id,category,severity,confidence,description,state,evidence,created_at,updated_at,tags)
        VALUES(?,?,?,?,?,?,?,?,?,?,?,?)
        ON CONFLICT(id) DO UPDATE SET severity=excluded.severity,confidence=excluded.confidence,
        description=excluded.description,state=excluded.state,evidence=excluded.evidence,updated_at=excluded.updated_at,tags=excluded.tags`,
		finding.ID, finding.SurveyID, finding.EchoID, finding.Category, finding.Severity,
		finding.Confidence, finding.Description, finding.State, evidence,
		timeText(finding.CreatedAt), timeText(finding.UpdatedAt), tags)
	if err != nil {
		return fmt.Errorf("save finding: %w", err)
	}
	return nil
}

func (r *Repository) GetFinding(ctx context.Context, id string) (model.Finding, error) {
	var finding model.Finding
	var evidence, tags, created, updated string
	err := r.db.QueryRowContext(ctx, `
        SELECT id,survey_id,echo_id,category,severity,confidence,description,state,evidence,created_at,updated_at,tags
        FROM findings WHERE id=?`, id).Scan(&finding.ID, &finding.SurveyID, &finding.EchoID, &finding.Category,
		&finding.Severity, &finding.Confidence, &finding.Description, &finding.State, &evidence, &created, &updated, &tags)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Finding{}, ErrNotFound
	}
	if err != nil {
		return model.Finding{}, fmt.Errorf("get finding: %w", err)
	}
	finding.CreatedAt = parseTime(created)
	finding.UpdatedAt = parseTime(updated)
	if err := unmarshal(evidence, &finding.Evidence); err != nil {
		return model.Finding{}, err
	}
	if err := unmarshal(tags, &finding.Tags); err != nil {
		return model.Finding{}, err
	}
	return finding, nil
}

func (r *Repository) ListFindings(ctx context.Context, surveyID string, minSeverity int) ([]model.Finding, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id,survey_id,echo_id,category,severity,confidence,description,state,evidence,created_at,updated_at,tags
        FROM findings WHERE survey_id=? AND severity>=? ORDER BY severity DESC,created_at`, surveyID, minSeverity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Finding, 0)
	for rows.Next() {
		var finding model.Finding
		var evidence, tags, created, updated string
		if err := rows.Scan(&finding.ID, &finding.SurveyID, &finding.EchoID, &finding.Category, &finding.Severity, &finding.Confidence, &finding.Description, &finding.State, &evidence, &created, &updated, &tags); err != nil {
			return nil, err
		}
		finding.CreatedAt = parseTime(created)
		finding.UpdatedAt = parseTime(updated)
		if err := unmarshal(evidence, &finding.Evidence); err != nil {
			return nil, err
		}
		if err := unmarshal(tags, &finding.Tags); err != nil {
			return nil, err
		}
		out = append(out, finding)
	}
	return out, rows.Err()
}

func (r *Repository) SaveReviewTrail(ctx context.Context, trail model.ReviewTrail) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO review_trails(id,finding_id,from_state,to_state,reviewer,reason,created_at,automatic,confidence) VALUES(?,?,?,?,?,?,?,?,?)`, trail.ID, trail.FindingID, trail.From, trail.To, trail.Reviewer, trail.Reason, timeText(trail.CreatedAt), boolInt(trail.Automatic), trail.Confidence)
	return err
}

func (r *Repository) CountOpenFindings(ctx context.Context, surveyID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM findings WHERE survey_id=? AND state=?", surveyID, model.AnomalyOpen).Scan(&count)
	return count, err
}
