package api

import "strconv"

type Page struct {
	Offset int
	Limit  int
	Total  int
}

func (p Page) Valid() bool    { return p.Offset >= 0 && p.Limit > 0 && p.Total >= 0 }
func (p Page) String() string { return strconv.Itoa(p.Offset) + ":" + strconv.Itoa(p.Limit) }
