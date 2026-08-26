package clock

import "time"

type Clock interface{ Now() time.Time }

type Real struct{}

func (Real) Now() time.Time { return time.Now().UTC() }

func StartOfDay(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}
