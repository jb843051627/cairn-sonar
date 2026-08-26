package model

import "strings"

type Operator struct {
	ID     string
	Name   string
	Scope  string
	Active bool
	Tags   []string
}

func (v Operator) Normalize() Operator {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v Operator) Ready() bool      { return v.ID != "" && v.Scope != "" && v.Active }
func (v Operator) Clone() Operator  { v.Tags = append([]string(nil), v.Tags...); return v }
func (v Operator) Describe() string { return v.ID + ":" + v.Scope }
