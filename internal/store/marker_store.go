package store

import (
	"context"
	"fmt"
)

func (r *Repository) PutMarker(ctx context.Context, id, scope string, active bool) error {
	if id == "" || scope == "" {
		return fmt.Errorf("markers requires id and scope")
	}
	_, err := r.db.ExecContext(ctx, "INSERT OR REPLACE INTO markers_markers(id, scope, active) VALUES(?,?,?)", id, scope, active)
	return err
}

func (r *Repository) CountMarker(ctx context.Context, scope string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM markers_markers WHERE scope=?", scope).Scan(&n)
	return n, err
}
