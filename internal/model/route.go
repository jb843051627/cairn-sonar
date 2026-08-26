package model

import "time"

type RouteStop struct {
	ChamberID string
	Order     int
	DurationM int
}

type Route struct {
	ID        string
	SurveyID  string
	Status    string
	Stops     []RouteStop
	DistanceM float64
	CreatedAt time.Time
}

func (r Route) Clone() Route {
	r.Stops = append([]RouteStop(nil), r.Stops...)
	return r
}

func (r Route) Valid() bool {
	return r.ID != "" && r.SurveyID != "" && len(r.Stops) > 0 && r.DistanceM >= 0
}
