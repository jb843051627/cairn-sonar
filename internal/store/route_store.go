package store

import (
	"cairn-sonar/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) SaveRoute(ctx context.Context, route model.Route) error {
	stops, err := marshal(route.Stops)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO routes(id,survey_id,status,stops,distance_m,created_at) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET status=excluded.status,stops=excluded.stops,distance_m=excluded.distance_m`, route.ID, route.SurveyID, route.Status, stops, route.DistanceM, timeText(route.CreatedAt))
	return err
}
func (r *Repository) GetRoute(ctx context.Context, id string) (model.Route, error) {
	var route model.Route
	var stops, created string
	err := r.db.QueryRowContext(ctx, `SELECT id,survey_id,status,stops,distance_m,created_at FROM routes WHERE id=?`, id).Scan(&route.ID, &route.SurveyID, &route.Status, &stops, &route.DistanceM, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Route{}, ErrNotFound
	}
	if err != nil {
		return model.Route{}, fmt.Errorf("get route: %w", err)
	}
	route.CreatedAt = parseTime(created)
	if err := unmarshal(stops, &route.Stops); err != nil {
		return model.Route{}, err
	}
	return route, nil
}
