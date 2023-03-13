package event

import (
	"context"
)

// Retry a function until either the returned error is nil or the context is done.
// The provided task is expected to be context aware.
func Retry(ctx context.Context, task func(context.Context) error) {
	for {
		if ctx.Err() != nil {
			return
		}
		// Task should be context aware, which means it should return
		// when ctx is cancelled.
		if err := task(ctx); err == nil {
			return
		}
	}
}
