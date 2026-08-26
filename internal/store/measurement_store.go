package store

import (
	"context"
	"fmt"
)

func (r *Repository) PutMeasurement(ctx context.Context, id, scope string, active bool) error {
	if id == "" || scope == "" {
		return fmt.Errorf("measurements requires id and scope")
	}
	_, err := r.db.ExecContext(ctx, "INSERT OR REPLACE INTO measurements_markers(id, scope, active) VALUES(?,?,?)", id, scope, active)
	return err
}

func (r *Repository) CountMeasurement(ctx context.Context, scope string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM measurements_markers WHERE scope=?", scope).Scan(&n)
	return n, err
}
