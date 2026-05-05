package job

import (
	"sync"
	"time"
)

// Status represents the current state of a monitored cron job.
type Status int

const (
	StatusUnknown Status = iota
	StatusOK
	StatusMissed
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusMissed:
		return "missed"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// Job holds the runtime state for a single monitored cron job.
type Job struct {
	mu          sync.Mutex
	Name        string
	Schedule    string
	GracePeriod time.Duration
	LastSeen    time.Time
	LastStatus  Status
	FailCount   int
}

// NewJob creates a Job from configuration values.
func NewJob(name, schedule string, gracePeriod time.Duration) *Job {
	return &Job{
		Name:        name,
		Schedule:    schedule,
		GracePeriod: gracePeriod,
		LastStatus:  StatusUnknown,
	}
}

// RecordSuccess marks the job as successfully executed at the given time.
func (j *Job) RecordSuccess(t time.Time) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.LastSeen = t
	j.LastStatus = StatusOK
	j.FailCount = 0
}

// RecordFailure increments the failure count and updates the status.
func (j *Job) RecordFailure(t time.Time) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.LastSeen = t
	j.LastStatus = StatusFailed
	j.FailCount++
}

// RecordMissed marks the job as having missed its scheduled run.
func (j *Job) RecordMissed() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.LastStatus = StatusMissed
	j.FailCount++
}

// Snapshot returns a copy of the current job state without holding the lock.
func (j *Job) Snapshot() Job {
	j.mu.Lock()
	defer j.mu.Unlock()
	return Job{
		Name:        j.Name,
		Schedule:    j.Schedule,
		GracePeriod: j.GracePeriod,
		LastSeen:    j.LastSeen,
		LastStatus:  j.LastStatus,
		FailCount:   j.FailCount,
	}
}
