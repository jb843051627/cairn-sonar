package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

var ErrNotFound = errors.New("record not found")

type Repository struct {
	db   *sql.DB
	path string
}

func Open(path string) (*Repository, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	repo := &Repository{db: db, path: path}
	if err := repo.configure(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	if err := repo.createSchema(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *Repository) configure(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000; PRAGMA journal_mode=WAL;")
	return err
}

func (r *Repository) Close() error { return r.db.Close() }
func (r *Repository) Path() string { return r.path }
func (r *Repository) DB() *sql.DB  { return r.db }

func (r *Repository) createSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}
