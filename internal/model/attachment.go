package model

import "strings"

type Attachment struct {
	ID     string
	Key    string
	Scope  string
	Active bool
	Tags   []string
}

func (v Attachment) Normalize() Attachment {
	v.ID = strings.TrimSpace(v.ID)
	v.Scope = strings.TrimSpace(v.Scope)
	v.Tags = append([]string(nil), v.Tags...)
	return v
}
func (v Attachment) Ready() bool       { return v.ID != "" && v.Scope != "" && v.Active }
func (v Attachment) Clone() Attachment { v.Tags = append([]string(nil), v.Tags...); return v }
func (v Attachment) Describe() string  { return v.ID + ":" + v.Scope }
