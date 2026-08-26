package model

import "strings"

type SiteMarker struct {
	ID     string
	Code   string
	Scope  string
	Active bool
	Tags   []string
}

func (v SiteMarker) Normalize() SiteMarker {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v SiteMarker) Ready() bool       { return v.ID != "" && v.Scope != "" && v.Active }
func (v SiteMarker) Clone() SiteMarker { v.Tags = append([]string(nil), v.Tags...); return v }
func (v SiteMarker) Describe() string  { return v.ID + ":" + v.Scope }
