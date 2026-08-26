package model

import "strings"

type WallSegment struct {
	ID      string
	LengthM float64
	Scope   string
	Active  bool
	Tags    []string
}

func (v WallSegment) Normalize() WallSegment {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v WallSegment) Ready() bool        { return v.ID != "" && v.Scope != "" && v.Active }
func (v WallSegment) Clone() WallSegment { v.Tags = append([]string(nil), v.Tags...); return v }
func (v WallSegment) Describe() string   { return v.ID + ":" + v.Scope }
