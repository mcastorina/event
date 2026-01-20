package event

import (
	"sync"
)

// Broadcast allows multiple goroutines to wait for a unique message. Once the
// message has been broadcasted, Wait() will always return.
type Broadcast struct {
	msgs map[string]struct{}
	cond *sync.Cond
}

// NewBroadcast initializes a [Broadcast] object.
func NewBroadcast() *Broadcast {
	return &Broadcast{
		msgs: make(map[string]struct{}),
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

// Broadcast announces the provided message and wakes any applicable goroutines
// waiting for the same message.
func (b *Broadcast) Broadcast(msg string) {
	// Acquire the lock to save the message.
	b.cond.L.Lock()
	defer b.cond.L.Unlock()
	b.msgs[msg] = struct{}{}
	// Broadcast that the map has changed.
	b.cond.Broadcast()
}

// Wait blocks until the provided message has been sent via [Broadcast]. If the
// message has already been broadcasted, then Wait returns immediately.
func (b *Broadcast) Wait(msg string) {
	for !b.waitForMsg(msg) {
	}
}

// waitForMsg is a helper method to check if a message was broadcasted and if
// not, block until the map of messages changes.
func (b *Broadcast) waitForMsg(msg string) bool {
	// Acquire the lock to check the map.
	b.cond.L.Lock()
	defer b.cond.L.Unlock()

	// If this message was already broadcasted, return.
	if _, ok := b.msgs[msg]; ok {
		return true
	}

	// Otherwise, wait for a signal that the map has changed.
	b.cond.Wait()

	// Return whether the message is in the map now or not.
	_, ok := b.msgs[msg]
	return ok
}
