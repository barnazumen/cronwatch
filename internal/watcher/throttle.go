// Package watcher contains the core monitoring loop and supporting types.
package watcher

import (
	"sync"
	"time"
)

// Throttle limits how frequently an action may fire for a given key.
// Unlike RateLimit (which uses a fixed cooldown window), Throttle enforces a
// minimum interval between consecutive firings and silently drops excess calls.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// NewThrottle creates a Throttle that allows at most one firing per interval
// for each distinct key.
func NewThrottle(interval time.Duration) *Throttle {
	return newThrottleWithClock(interval, time.Now)
}

func newThrottleWithClock(interval time.Duration, now func() time.Time) *Throttle {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Throttle{
		interval: interval,
		last:     make(map[string]time.Time),
		now:      now,
	}
}

// Allow returns true if the key has not fired within the throttle interval.
// When true is returned the internal timestamp for the key is updated.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok && now.Sub(last) < t.interval {
		return false
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded timestamp for key, allowing the next call to
// Allow to succeed immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}
