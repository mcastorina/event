package event

import (
	"sync"
)

// Trigger allows multiple goroutines to synchronize on a single event. Once
// the event has been triggered, there is no resetting it.
type Trigger struct {
	initLock  sync.Mutex
	cond      *sync.Cond
	triggered bool
}

// Trigger triggers the event, signalling to all waiting goroutines to
// continue.
func (t *Trigger) Trigger() {
	if t.triggered {
		return
	}
	t.ensureInit()
	t.cond.L.Lock()
	t.triggered = true
	t.cond.L.Unlock()
	t.cond.Broadcast()
}

// Wait for the event to trigger. If the event has already been triggered, Wait
// returns immediately.
func (t *Trigger) Wait() {
	if t.triggered {
		return
	}
	t.ensureInit()
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	for !t.triggered {
		t.cond.Wait()
	}
}

// Done is a convenience method for waiting via a channel.
func (t *Trigger) Done() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		t.Wait()
		ch <- struct{}{}
	}()
	return ch
}

// ensureInit is a helper method to ensure the internal sync.Cond is
// initialized.
func (t *Trigger) ensureInit() {
	if t.cond != nil {
		return
	}
	t.initLock.Lock()
	defer t.initLock.Unlock()
	if t.cond != nil {
		return
	}
	t.cond = sync.NewCond(&sync.Mutex{})
}
