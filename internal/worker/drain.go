package worker

import (
	"context"
	"time"
)

func DrainTick(ctx context.Context, interval time.Duration, fn func()) error {
	if interval <= 0 {
		interval = time.Second
	}
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		fn()
		return nil
	}
}

func DrainLabel() string { return "drain" }
