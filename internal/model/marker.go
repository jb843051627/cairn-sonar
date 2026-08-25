package model

import "strings"

type WallMarker struct {
	ID     string
	Scope  string
	Active bool
	Tags   []string
}

func (v WallMarker) Normalize() WallMarker {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v WallMarker) Ready() bool       { return v.ID != "" && v.Scope != "" && v.Active }
func (v WallMarker) Clone() WallMarker { v.Tags = append([]string(nil), v.Tags...); return v }
func (v WallMarker) Describe() string  { return v.ID + ":" + v.Scope }
