package rules

import (
	"cairn-sonar/internal/model"
	"sort"
)

func BuildRoute(chambers []model.Chamber) model.Route {
	ordered := append([]model.Chamber(nil), chambers...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].DepthM < ordered[j].DepthM })
	stops := make([]model.RouteStop, 0, len(ordered))
	distance := 0.0
	for i, chamber := range ordered {
		stops = append(stops, model.RouteStop{ChamberID: chamber.ID, Order: i + 1, DurationM: 15})
		distance += chamber.DepthM + 2
	}
	return model.Route{Stops: stops, DistanceM: distance, Status: "planned"}
}
