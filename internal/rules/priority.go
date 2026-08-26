package rules

import "sort"

func PriorityOrder(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}

func PriorityRank(value int) string {
	switch {
	case value >= 80:
		return "high"
	case value >= 50:
		return "medium"
	default:
		return "low"
	}
}
