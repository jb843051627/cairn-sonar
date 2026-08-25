package api

import (
	"fmt"
	"net/url"
)

type Query struct{ Values url.Values }

func ParseQuery(values url.Values) Query { return Query{Values: values} }

func (q Query) String(name, fallback string) string {
	value := q.Values.Get(name)
	if value == "" {
		return fallback
	}
	return value
}

func (q Query) Int(name string, fallback int) int {
	value := q.Values.Get(name)
	var n int
	if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
		return fallback
	}
	return n
}
