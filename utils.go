package event

import (
	"context"
	"errors"
	"time"
)

// Poll asynchronously executes predicate at the given period until true is
// returned.
func Poll(period time.Duration, predicate func() bool) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for {
			if predicate() {
				ch <- struct{}{}
				return
			}
			time.Sleep(period)
		}
	}()
	return ch
}

// Timeout executes task and returns an error on timeout. If the task executes
// before the timeout, nil is returned. The context passed to the task will be
// canceled on timeout as well.
func Timeout(wait time.Duration, task func(context.Context)) error {
	if task == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	go func() {
		task(ctx)
		cancel()
	}()
	<-ctx.Done()
	if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return nil
}
