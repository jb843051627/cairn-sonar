package rules

import (
	"math"
	"sort"

	"cairn-sonar/internal/model"
)

type RoutePolicy struct {
	MaxStops       int
	MaxDistanceM   float64
	PreferDeep     bool
	IncludeClosed  bool
	ReturnToOrigin bool
}

func DefaultRoutePolicy() RoutePolicy {
	return RoutePolicy{MaxStops: 24, MaxDistanceM: 2000, PreferDeep: true, ReturnToOrigin: true}
}

type RouteCandidate struct {
	Chamber model.Chamber
	Score   float64
	Reason  string
}

func RankChambers(chambers []model.Chamber, policy RoutePolicy) []RouteCandidate {
	out := make([]RouteCandidate, 0, len(chambers))
	for _, chamber := range chambers {
		if !chamber.Valid() {
			continue
		}
		score := chamber.DepthM
		reason := "depth priority"
		if !policy.PreferDeep {
			score = -chamber.DepthM
			reason = "shallow priority"
		}
		if chamber.Temperature > 28 {
			score += 10
			reason += ", thermal note"
		}
		out = append(out, RouteCandidate{Chamber: chamber.Clone(), Score: score, Reason: reason})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func BuildPolicyRoute(chambers []model.Chamber, policy RoutePolicy) model.Route {
	candidates := RankChambers(chambers, policy)
	if policy.MaxStops <= 0 {
		policy.MaxStops = len(candidates)
	}
	stops := make([]model.RouteStop, 0, min(policy.MaxStops, len(candidates)))
	distance := 0.0
	for i, candidate := range candidates {
		if len(stops) >= policy.MaxStops {
			break
		}
		increment := math.Abs(candidate.Chamber.DepthM) + 2
		if policy.ReturnToOrigin && i > 0 {
			increment += math.Abs(candidate.Chamber.DepthM-candidates[i-1].Chamber.DepthM) / 2
		}
		if policy.MaxDistanceM > 0 && distance+increment > policy.MaxDistanceM {
			continue
		}
		distance += increment
		stops = append(stops, model.RouteStop{ChamberID: candidate.Chamber.ID, Order: len(stops) + 1, DurationM: stopDuration(candidate.Chamber)})
	}
	return model.Route{Stops: stops, DistanceM: distance, Status: "planned"}
}

func stopDuration(chamber model.Chamber) int {
	duration := 15
	if chamber.DepthM > 40 {
		duration += 10
	}
	if chamber.Temperature > 28 {
		duration += 5
	}
	return duration
}

func EstimateRouteMinutes(route model.Route, travelSpeedMPerMin float64) int {
	if travelSpeedMPerMin <= 0 {
		travelSpeedMPerMin = 20
	}
	minutes := route.DistanceM / travelSpeedMPerMin
	for _, stop := range route.Stops {
		minutes += float64(stop.DurationM)
	}
	return int(math.Ceil(minutes))
}

func ValidateRoute(route model.Route, policy RoutePolicy) []string {
	errors := make([]string, 0)
	if len(route.Stops) == 0 {
		errors = append(errors, "route has no stops")
	}
	if policy.MaxStops > 0 && len(route.Stops) > policy.MaxStops {
		errors = append(errors, "route exceeds stop limit")
	}
	if policy.MaxDistanceM > 0 && route.DistanceM > policy.MaxDistanceM {
		errors = append(errors, "route exceeds distance limit")
	}
	previous := 0
	for _, stop := range route.Stops {
		if stop.Order != previous+1 {
			errors = append(errors, "route order is not contiguous")
		}
		if stop.DurationM <= 0 {
			errors = append(errors, "route stop has no duration")
		}
		previous = stop.Order
	}
	return errors
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
