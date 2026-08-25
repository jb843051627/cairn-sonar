package store

import (
	"context"
	"fmt"
)

func (r *Repository) PutSummary(ctx context.Context, id, scope string, active bool) error {
	if id == "" || scope == "" {
		return fmt.Errorf("summaries requires id and scope")
	}
	_, err := r.db.ExecContext(ctx, "INSERT OR REPLACE INTO summaries_markers(id, scope, active) VALUES(?,?,?)", id, scope, active)
	return err
}

func (r *Repository) CountSummary(ctx context.Context, scope string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM summaries_markers WHERE scope=?", scope).Scan(&n)
	return n, err
}
