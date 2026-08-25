package worker

import "sync"

type Barrier struct {
	mu   sync.Mutex
	wait sync.WaitGroup
}

func (b *Barrier) Add(n int) { b.mu.Lock(); b.wait.Add(n); b.mu.Unlock() }
func (b *Barrier) Done()     { b.wait.Done() }
func (b *Barrier) Wait()     { b.wait.Wait() }
