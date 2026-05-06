package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/cronwatch/internal/config"
)

// Status represents the last known state of a job.
type Status int

const (
	StatusOK      Status = iota
	StatusFailed
	StatusMissed
	StatusUnknown
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusFailed:
		return "failed"
	case StatusMissed:
		return "missed"
	default:
		return "unknown"
	}
}

// Job holds runtime state for a monitored cron job.
type Job struct {
	mu          sync.Mutex
	name        string
	maxInterval time.Duration
	lastRun     time.Time
	status      Status
	failCount   int
}

// NewJob constructs a Job from a JobConfig.
func NewJob(cfg config.JobConfig) (*Job, error) {
	d, err := time.ParseDuration(cfg.MaxInterval)
	if err != nil {
		return nil, fmt.Errorf("job %q: invalid max_interval %q: %w", cfg.Name, cfg.MaxInterval, err)
	}
	return &Job{
		name:        cfg.Name,
		maxInterval: d,
		lastRun:     time.Now(),
		status:      StatusUnknown,
	}, nil
}

func (j *Job) Name() string        { return j.name }
func (j *Job) MaxInterval() time.Duration {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.maxInterval
}
func (j *Job) LastRun() time.Time {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.lastRun
}
func (j *Job) Status() Status {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.status
}
func (j *Job) FailCount() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.failCount
}

// RecordSuccess marks the job as healthy and updates the last-run timestamp.
func (j *Job) RecordSuccess() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.lastRun = time.Now()
	j.status = StatusOK
	j.failCount = 0
}

// RecordFailure increments the failure counter.
func (j *Job) RecordFailure() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.lastRun = time.Now()
	j.status = StatusFailed
	j.failCount++
}

// RecordMissed marks the job as having missed its scheduled window.
func (j *Job) RecordMissed() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.status = StatusMissed
}
