package watcher

import (
	"sync"
	"time"
)

// RateLimit suppresses repeated alerts for the same job within a cooldown
// window, preventing alert storms when a job fails on every check cycle.
type RateLimit struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
	now      func() time.Time
}

// NewRateLimit creates a RateLimit with the given cooldown duration.
// Alerts for a given job name are suppressed until cooldown has elapsed
// since the last alert was forwarded.
func NewRateLimit(cooldown time.Duration) *RateLimit {
	return &RateLimit{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if an alert for jobName should be forwarded.
// It records the current time as the last-sent time when it returns true.
func (r *RateLimit) Allow(jobName string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	last, seen := r.lastSent[jobName]
	if !seen || r.now().Sub(last) >= r.cooldown {
		r.lastSent[jobName] = r.now()
		return true
	}
	return false
}

// Reset clears the rate-limit state for jobName, allowing the next alert
// through immediately. Useful when a job recovers.
func (r *RateLimit) Reset(jobName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.lastSent, jobName)
}

// Snapshot returns a copy of the last-sent map for inspection or testing.
func (r *RateLimit) Snapshot() map[string]time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]time.Time, len(r.lastSent))
	for k, v := range r.lastSent {
		out[k] = v
	}
	return out
}
