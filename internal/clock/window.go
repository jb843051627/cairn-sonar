package clock

import "time"

type Window struct{ From, To time.Time }

func NewWindow(from, to time.Time) Window  { return Window{From: from.UTC(), To: to.UTC()} }
func (w Window) Contains(t time.Time) bool { u := t.UTC(); return !u.Before(w.From) && u.Before(w.To) }
func (w Window) Duration() time.Duration   { return w.To.Sub(w.From) }
