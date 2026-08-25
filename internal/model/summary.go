package model

import "strings"

type SurveySummary struct {
	ID     string
	Count  int
	Scope  string
	Active bool
	Tags   []string
}

func (v SurveySummary) Normalize() SurveySummary {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v SurveySummary) Ready() bool          { return v.ID != "" && v.Scope != "" && v.Active }
func (v SurveySummary) Clone() SurveySummary { v.Tags = append([]string(nil), v.Tags...); return v }
func (v SurveySummary) Describe() string     { return v.ID + ":" + v.Scope }
