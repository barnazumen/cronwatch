// Package watcher contains the core monitoring loop and supporting types.
package watcher

import (
	"sync"
	"time"
)

// Dedup suppresses duplicate alerts for the same job within a fixed window.
// Unlike RateLimit (which uses a cooldown per job), Dedup tracks the last
// alert content hash so that identical back-to-back alerts are silently
// dropped even if the cooldown has expired.
type Dedup struct {
	mu      sync.Mutex
	window  time.Duration
	records map[string]dedupRecord
}

type dedupRecord struct {
	key     string
	seenAt  time.Time
}

// NewDedup returns a Dedup that suppresses identical alerts within window.
func NewDedup(window time.Duration) *Dedup {
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &Dedup{
		window:  window,
		records: make(map[string]dedupRecord),
	}
}

// IsDuplicate returns true when an alert with the same jobName+key has been
// seen within the dedup window. It records the alert if it is not a duplicate.
func (d *Dedup) IsDuplicate(jobName, key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	composite := jobName + "\x00" + key
	now := time.Now()

	if rec, ok := d.records[composite]; ok {
		if now.Sub(rec.seenAt) < d.window {
			return true
		}
	}

	d.records[composite] = dedupRecord{key: key, seenAt: now}
	return false
}

// Reset clears the dedup state for a specific job, allowing the next alert
// through regardless of the window.
func (d *Dedup) Reset(jobName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for k := range d.records {
		// keys are "jobName\x00key"
		if len(k) > len(jobName) && k[:len(jobName)] == jobName && k[len(jobName)] == 0 {
			delete(d.records, k)
		}
	}
}
