// Package watcher provides the core loop that checks cron jobs for missed runs.
package watcher

import (
	"sync"
	"time"
)

// CheckRecord holds the result of a single watcher check cycle.
type CheckRecord struct {
	At      time.Time
	Missed  []string
	Alerted []string
}

// History stores a bounded ring-buffer of recent check records.
type History struct {
	mu      sync.Mutex
	records []CheckRecord
	cap     int
}

// NewHistory creates a History that retains at most maxRecords entries.
func NewHistory(maxRecords int) *History {
	if maxRecords <= 0 {
		maxRecords = 100
	}
	return &History{
		records: make([]CheckRecord, 0, maxRecords),
		cap:     maxRecords,
	}
}

// Record appends a CheckRecord, evicting the oldest entry when full.
func (h *History) Record(r CheckRecord) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.records) >= h.cap {
		h.records = h.records[1:]
	}
	h.records = append(h.records, r)
}

// Snapshot returns a copy of all stored records, oldest first.
func (h *History) Snapshot() []CheckRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]CheckRecord, len(h.records))
	copy(out, h.records)
	return out
}

// Len returns the current number of stored records.
func (h *History) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.records)
}
