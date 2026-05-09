package watcher

import (
	"sync"
	"time"
)

// HealthStatus represents the current health state of a monitored job.
type HealthStatus struct {
	JobName   string
	Healthy   bool
	LastCheck time.Time
	Message   string
}

// HealthCheck tracks per-job health state based on recent check outcomes.
type HealthCheck struct {
	mu       sync.RWMutex
	states   map[string]*HealthStatus
	nowFn    func() time.Time
}

// NewHealthCheck creates a new HealthCheck instance.
func NewHealthCheck() *HealthCheck {
	return &HealthCheck{
		states: make(map[string]*HealthStatus),
		nowFn:  time.Now,
	}
}

// RecordHealthy marks a job as healthy.
func (h *HealthCheck) RecordHealthy(jobName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.states[jobName] = &HealthStatus{
		JobName:   jobName,
		Healthy:   true,
		LastCheck: h.nowFn(),
		Message:   "job ran on time",
	}
}

// RecordUnhealthy marks a job as unhealthy with a reason message.
func (h *HealthCheck) RecordUnhealthy(jobName, reason string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.states[jobName] = &HealthStatus{
		JobName:   jobName,
		Healthy:   false,
		LastCheck: h.nowFn(),
		Message:   reason,
	}
}

// Status returns the current health status for a job.
// The second return value is false if the job has never been recorded.
func (h *HealthCheck) Status(jobName string) (HealthStatus, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	s, ok := h.states[jobName]
	if !ok {
		return HealthStatus{}, false
	}
	return *s, true
}

// Snapshot returns a copy of all current health statuses.
func (h *HealthCheck) Snapshot() []HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]HealthStatus, 0, len(h.states))
	for _, s := range h.states {
		out = append(out, *s)
	}
	return out
}
