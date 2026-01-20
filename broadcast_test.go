package event

import (
	"sync"
	"testing"
)

func TestBroadcast(t *testing.T) {
	b := NewBroadcast()
	signalCounts := make([]int, 3)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		b.Wait("first")
		signalCounts[0]++
	}()
	go func() {
		defer wg.Done()
		b.Wait("first")
		signalCounts[1]++
		b.Wait("second")
		signalCounts[1]++
	}()
	go func() {
		defer wg.Done()
		b.Wait("third")
		signalCounts[2]++
		b.Wait("first")
		signalCounts[2]++
		b.Wait("second")
		signalCounts[2]++
	}()

	b.Broadcast("first")
	b.Broadcast("second")
	b.Broadcast("third")

	wg.Wait()

	if signalCounts[0] != 1 {
		t.Errorf("expected 1 signal, found %d", signalCounts[0])
	}
	if signalCounts[1] != 2 {
		t.Errorf("expected 2 signals, found %d", signalCounts[1])
	}
	if signalCounts[2] != 3 {
		t.Errorf("expected 3 signal, found %d", signalCounts[2])
	}
}
