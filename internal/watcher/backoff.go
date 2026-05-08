// Package watcher contains the core polling loop and supporting types.
package watcher

import "time"

// Backoff computes exponential back-off delays for repeated alert
// suppression so that a persistently-missed job does not flood the
// configured alerter.
//
// The delay starts at Base and doubles on every call to Next until it
// reaches Cap.  Reset returns the delay to Base.
type Backoff struct {
	Base    time.Duration
	Cap     time.Duration
	current time.Duration
}

// NewBackoff returns a Backoff with sensible defaults (1 min base, 1 h cap).
func NewBackoff(base, cap time.Duration) *Backoff {
	if base <= 0 {
		base = time.Minute
	}
	if cap < base {
		cap = base
	}
	return &Backoff{Base: base, Cap: cap, current: base}
}

// Current returns the delay that would be returned by the next call to Next
// without advancing the sequence.
func (b *Backoff) Current() time.Duration {
	return b.current
}

// Next returns the current delay and advances to the next (doubled) value.
func (b *Backoff) Next() time.Duration {
	d := b.current
	next := b.current * 2
	if next > b.Cap {
		next = b.Cap
	}
	b.current = next
	return d
}

// Reset returns the delay sequence to Base.
func (b *Backoff) Reset() {
	b.current = b.Base
}
