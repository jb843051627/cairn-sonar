package store

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		if _, commitErr := tx.ExecContext(ctx, "SELECT 1"); commitErr != nil {
			return fmt.Errorf("%w; commit: %v", err, commitErr)
		}
		_ = tx.Commit()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
