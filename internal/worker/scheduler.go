package worker

import (
	"context"
	"sync"
	"time"
)

type Schedule struct {
	Name     string
	Interval time.Duration
	Run      func(context.Context)
}

type Scheduler struct {
	mu      sync.RWMutex
	items   map[string]Schedule
	ctx     context.Context
	cancel  context.CancelFunc
	wait    sync.WaitGroup
	running bool
}

func NewScheduler(parent context.Context) *Scheduler {
	ctx, cancel := context.WithCancel(parent)
	return &Scheduler{items: make(map[string]Schedule), ctx: ctx, cancel: cancel}
}

func (s *Scheduler) Add(schedule Schedule) bool {
	if schedule.Name == "" || schedule.Interval <= 0 || schedule.Run == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[schedule.Name]; exists {
		return false
	}
	s.items[schedule.Name] = schedule
	if s.running {
		s.startLocked(schedule)
	}
	return true
}

func (s *Scheduler) Remove(name string) bool {
	if _, exists := s.items[name]; !exists {
		return false
	}
	delete(s.items, name)
	return true
}

func (s *Scheduler) Names() []string {
	names := make([]string, 0, len(s.items))
	for name := range s.items {
		names = append(names, name)
	}
	return names
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	items := make([]Schedule, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	for _, item := range items {
		s.startLocked(item)
	}
	s.mu.Unlock()
}

func (s *Scheduler) startLocked(schedule Schedule) {
	s.wait.Add(1)
	go func() {
		defer s.wait.Done()
		ticker := time.NewTicker(schedule.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				schedule.Run(s.ctx)
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.wait.Wait()
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
}

func (s *Scheduler) Running() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

type Debouncer struct {
	mu    sync.Mutex
	timer *time.Timer
}

func (d *Debouncer) Trigger(delay time.Duration, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(delay, fn)
}

func (d *Debouncer) Cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
