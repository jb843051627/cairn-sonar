package worker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrRetryExhausted = errors.New("retry exhausted")

func Retry(ctx context.Context, attempts int, delay time.Duration, fn func(context.Context) error) error {
	if attempts < 1 {
		attempts = 1
	}
	var last error
	for i := 0; i < attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(ctx); err == nil {
			return nil
		} else {
			last = err
		}
		if i+1 < attempts {
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return fmt.Errorf("%w: %v", ErrRetryExhausted, last)
}
