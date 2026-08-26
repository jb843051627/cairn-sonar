package model

import (
	"strings"
	"time"
)

type Chamber struct {
	ID          string
	Name        string
	SiteCode    string
	DepthM      float64
	Temperature float64
	CreatedAt   time.Time
	Tags        []string
}

func (c Chamber) Normalize() Chamber {
	c.ID = strings.TrimSpace(c.ID)
	c.Name = strings.TrimSpace(c.Name)
	c.SiteCode = strings.ToUpper(strings.TrimSpace(c.SiteCode))
	c.Tags = cloneStrings(c.Tags)
	return c
}

func (c Chamber) Clone() Chamber {
	c.Tags = cloneStrings(c.Tags)
	return c
}

func (c Chamber) Valid() bool {
	return c.ID != "" && c.Name != "" && c.SiteCode != "" && c.DepthM >= 0
}
