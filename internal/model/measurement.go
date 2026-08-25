package model

import "strings"

type Measurement struct {
	ID     string
	Value  float64
	Scope  string
	Active bool
	Tags   []string
}

func (v Measurement) Normalize() Measurement {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v Measurement) Ready() bool        { return v.ID != "" && v.Scope != "" && v.Active }
func (v Measurement) Clone() Measurement { v.Tags = append([]string(nil), v.Tags...); return v }
func (v Measurement) Describe() string   { return v.ID + ":" + v.Scope }
