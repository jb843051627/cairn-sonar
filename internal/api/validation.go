package api

import "strconv"

type RequestCheck struct {
	Offset int
	Limit  int
	Total  int
}

func (p RequestCheck) Valid() bool    { return p.Offset >= 0 && p.Limit > 0 && p.Total >= 0 }
func (p RequestCheck) String() string { return strconv.Itoa(p.Offset) + ":" + strconv.Itoa(p.Limit) }
