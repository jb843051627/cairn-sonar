package model

import "strings"

type FieldReading struct {
	ID     string
	Value  float64
	Scope  string
	Active bool
	Tags   []string
}

func (v FieldReading) Normalize() FieldReading {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v FieldReading) Ready() bool         { return v.ID != "" && v.Scope != "" && v.Active }
func (v FieldReading) Clone() FieldReading { v.Tags = append([]string(nil), v.Tags...); return v }
func (v FieldReading) Describe() string    { return v.ID + ":" + v.Scope }
