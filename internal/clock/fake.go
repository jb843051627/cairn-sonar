package clock

import (
	"sync"
	"time"
)

type Fake struct {
	mu      sync.RWMutex
	current time.Time
}

func NewFake(t time.Time) *Fake         { return &Fake{current: t.UTC()} }
func (f *Fake) Now() time.Time          { f.mu.RLock(); defer f.mu.RUnlock(); return f.current }
func (f *Fake) Advance(d time.Duration) { f.mu.Lock(); f.current = f.current.Add(d); f.mu.Unlock() }
func (f *Fake) Set(t time.Time)         { f.mu.Lock(); f.current = t.UTC(); f.mu.Unlock() }
