package event

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTrigger(t *testing.T) {
	var tr Trigger
	var wg sync.WaitGroup

	ready := make([]atomic.Bool, 4)
	done := make([]atomic.Bool, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ready[i].Store(true)
			tr.Wait()
			done[i].Store(true)
		}(i)
	}

	select {
	case <-time.After(1 * time.Second):
		t.Errorf("timeout waiting for all goroutines to be ready")
	case <-waitAllTrue(ready):
	}

	if !all(done, false) {
		t.Errorf("expected all goroutines to be waiting")
	}

	tr.Trigger()
	wg.Wait()

	if !all(done, true) {
		t.Errorf("expected all goroutines to be done")
	}
}

func all(arr []atomic.Bool, expected bool) bool {
	for i := range arr {
		if arr[i].Load() != expected {
			return false
		}
	}
	return true
}

func waitAllTrue(arr []atomic.Bool) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for !all(arr, true) {
		}
		ch <- struct{}{}
	}()
	return ch
}
