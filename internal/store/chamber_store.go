package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) PutChamber(ctx context.Context, c model.Chamber) error {
	tags, err := marshal(c.Tags)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO chambers(id,name,site_code,depth_m,temperature,created_at,tags) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,site_code=excluded.site_code,depth_m=excluded.depth_m,temperature=excluded.temperature,tags=excluded.tags`, c.ID, c.Name, c.SiteCode, c.DepthM, c.Temperature, timeText(c.CreatedAt), tags)
	return err
}

func (r *Repository) GetChamber(ctx context.Context, id string) (model.Chamber, error) {
	var c model.Chamber
	var created, tags string
	err := r.db.QueryRowContext(ctx, `SELECT id,name,site_code,depth_m,temperature,created_at,tags FROM chambers WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.SiteCode, &c.DepthM, &c.Temperature, &created, &tags)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Chamber{}, ErrNotFound
	}
	if err != nil {
		return model.Chamber{}, fmt.Errorf("get chamber: %w", err)
	}
	c.CreatedAt = parseTime(created)
	if err := unmarshal(tags, &c.Tags); err != nil {
		return model.Chamber{}, err
	}
	return c, nil
}

func (r *Repository) ListChambers(ctx context.Context) ([]model.Chamber, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,name,site_code,depth_m,temperature,created_at,tags FROM chambers ORDER BY depth_m,id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Chamber, 0)
	for rows.Next() {
		var c model.Chamber
		var created, tags string
		if err := rows.Scan(&c.ID, &c.Name, &c.SiteCode, &c.DepthM, &c.Temperature, &created, &tags); err != nil {
			return nil, err
		}
		c.CreatedAt = parseTime(created)
		if err := unmarshal(tags, &c.Tags); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
