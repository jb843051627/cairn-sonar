package model

import "strings"

type TaxonomyLabel struct {
	ID     string
	Value  string
	Scope  string
	Active bool
	Tags   []string
}

func (v TaxonomyLabel) Normalize() TaxonomyLabel {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v TaxonomyLabel) Ready() bool          { return v.ID != "" && v.Scope != "" && v.Active }
func (v TaxonomyLabel) Clone() TaxonomyLabel { v.Tags = append([]string(nil), v.Tags...); return v }
func (v TaxonomyLabel) Describe() string     { return v.ID + ":" + v.Scope }
