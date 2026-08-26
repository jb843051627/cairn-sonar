package worker

import (
	"context"
	"errors"
	"sync"
)

var ErrBatchCancelled = errors.New("batch cancelled")

type BatchResult[T any] struct {
	Index int
	Value T
	Err   error
}

func Map[T any, R any](ctx context.Context, values []T, workers int, fn func(context.Context, T) (R, error)) []BatchResult[R] {
	if workers < 1 {
		workers = 1
	}
	if workers > len(values) && len(values) > 0 {
		workers = len(values)
	}
	results := make([]BatchResult[R], len(values))
	jobs := make(chan int)
	var wait sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for index := range jobs {
				if err := ctx.Err(); err != nil {
					results[index] = BatchResult[R]{Index: index, Err: err}
					continue
				}
				value, err := fn(ctx, values[index])
				results[index] = BatchResult[R]{Index: index, Value: value, Err: err}
			}
		}()
	}
	for index := range values {
		select {
		case <-ctx.Done():
			results[index] = BatchResult[R]{Index: index, Err: ErrBatchCancelled}
		case jobs <- index:
		}
	}
	close(jobs)
	wait.Wait()
	return results
}

type Gate struct {
	mu     sync.Mutex
	closed bool
	items  int
}

func (g *Gate) Acquire() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return false
	}
	g.items++
	return true
}

func (g *Gate) Release() {
	g.mu.Lock()
	if g.items > 0 {
		g.items--
	}
	g.mu.Unlock()
}

func (g *Gate) Close() {
	g.mu.Lock()
	g.closed = true
	g.mu.Unlock()
}

func (g *Gate) Active() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.items
}
