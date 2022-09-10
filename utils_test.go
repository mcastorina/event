package event

import (
	"context"
	"testing"
	"time"
)

func TestPoll(t *testing.T) {
	var i int
	<-Poll(1*time.Nanosecond, func() bool {
		i++
		return i >= 3
	})
	if i != 3 {
		t.Errorf("expected to poll 3 times")
	}
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		task        func(context.Context)
		wantTimeout bool
	}{
		{
			name:        "empty task",
			timeout:     1 * time.Second,
			task:        func(_ context.Context) {},
			wantTimeout: false,
		},
		{
			name:        "nil task",
			timeout:     1 * time.Second,
			task:        nil,
			wantTimeout: false,
		},
		{
			name:        "timeout",
			timeout:     0,
			task:        func(_ context.Context) { time.Sleep(1 * time.Second) },
			wantTimeout: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Timeout(tt.timeout, tt.task)
			haveTimeout := err != nil
			if haveTimeout != tt.wantTimeout {
				t.Fail()
			}
		})
	}
}
