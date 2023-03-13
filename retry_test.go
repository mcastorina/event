package event

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	var count int
	Retry(context.Background(), func(context.Context) error {
		count++
		if count < 3 {
			return fmt.Errorf("oh no")
		}
		return nil
	})

	if count != 3 {
		t.Errorf("expected count to be 3, got %d", count)
	}
}

func TestRetryTimeout(t *testing.T) {
	var count int
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	Retry(ctx, func(context.Context) error {
		count++
		time.Sleep(50 * time.Millisecond)
		return fmt.Errorf("oh no")
	})
	if count != 2 {
		t.Errorf("expected count to be 2, got %d", count)
	}
}
