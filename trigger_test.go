package event

import (
	"sync"
	"testing"
	"time"
)

func TestTrigger(t *testing.T) {
	var tr Trigger
	var wg sync.WaitGroup

	ready := make([]bool, 4)
	done := make([]bool, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ready[i] = true
			tr.Wait()
			done[i] = true
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

func all(arr []bool, expected bool) bool {
	for _, b := range arr {
		if b != expected {
			return false
		}
	}
	return true
}

func waitAllTrue(arr []bool) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for !all(arr, true) {
		}
		ch <- struct{}{}
	}()
	return ch
}
