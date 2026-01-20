package event

import (
	"context"
	"sync"
)

// Trigger allows multiple goroutines to synchronize on a single event. Once
// the event has been triggered, there is no resetting it.
type Trigger struct {
	initLock sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// Trigger triggers the event, signalling to all waiting goroutines to
// continue.
func (t *Trigger) Trigger() {
	t.ensureInit()
	t.cancel()
}

// Wait for the event to trigger. If the event has already been triggered, Wait
// returns immediately.
func (t *Trigger) Wait() {
	t.ensureInit()
	<-t.ctx.Done()
}

// Done is a convenience method for waiting via a channel.
func (t *Trigger) Done() <-chan struct{} {
	t.ensureInit()
	return t.ctx.Done()
}

// ensureInit is a helper method to ensure the internal sync.Cond is
// initialized.
func (t *Trigger) ensureInit() {
	t.initLock.Lock()
	defer t.initLock.Unlock()
	if t.ctx != nil {
		return
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
}
