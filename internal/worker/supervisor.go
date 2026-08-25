package worker

import (
	"context"
	"sync"
)

type Supervisor struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSupervisor(parent context.Context) (*Supervisor, context.Context) {
	ctx, cancel := context.WithCancel(parent)
	return &Supervisor{ctx: ctx, cancel: cancel}, ctx
}
func (s *Supervisor) Go(fn func(context.Context)) {
	s.wg.Add(1)
	go func() { defer s.wg.Done(); fn(s.ctx) }()
}
func (s *Supervisor) Stop() { s.cancel(); s.wg.Wait() }
