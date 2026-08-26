package model

import "strings"

type Checkpoint struct {
	ID     string
	Label  string
	Scope  string
	Active bool
	Tags   []string
}

func (v Checkpoint) Normalize() Checkpoint {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v Checkpoint) Ready() bool       { return v.ID != "" && v.Scope != "" && v.Active }
func (v Checkpoint) Clone() Checkpoint { v.Tags = append([]string(nil), v.Tags...); return v }
func (v Checkpoint) Describe() string  { return v.ID + ":" + v.Scope }
