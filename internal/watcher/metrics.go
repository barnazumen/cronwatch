package watcher

import "sync/atomic"

// Metrics tracks runtime counters for the watcher.
type Metrics struct {
	Checks  atomic.Int64
	Missed  atomic.Int64
	Alerted atomic.Int64
}

// NewMetrics returns a zeroed Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// RecordCheck increments the total check counter.
func (m *Metrics) RecordCheck() {
	m.Checks.Add(1)
}

// RecordMissed increments the missed-run counter.
func (m *Metrics) RecordMissed() {
	m.Missed.Add(1)
}

// RecordAlerted increments the alert-sent counter.
func (m *Metrics) RecordAlerted() {
	m.Alerted.Add(1)
}

// Snapshot returns a point-in-time copy of the counters.
type MetricsSnapshot struct {
	Checks  int64
	Missed  int64
	Alerted int64
}

// Snapshot reads all counters atomically.
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		Checks:  m.Checks.Load(),
		Missed:  m.Missed.Load(),
		Alerted: m.Alerted.Load(),
	}
}
