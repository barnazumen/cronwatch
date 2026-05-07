// Package watcher provides the core loop that checks job schedules
// and fires alerts when runs are missed or fail.
package watcher

import (
	"fmt"
	"time"

	"github.com/cronwatch/cronwatch/internal/job"
)

// EventKind classifies what happened during a watcher cycle.
type EventKind string

const (
	// EventMissed means the job did not run within its expected window.
	EventMissed EventKind = "missed"
	// EventFailure means the job ran but reported a non-zero exit.
	EventFailure EventKind = "failure"
	// EventRecovered means the job is healthy again after a prior alert.
	EventRecovered EventKind = "recovered"
)

// Event carries information about a single watcher observation.
type Event struct {
	// Kind is the classification of this event.
	Kind EventKind
	// Job is the job that triggered the event.
	Job *job.Job
	// ObservedAt is when the watcher detected the condition.
	ObservedAt time.Time
	// Message is a human-readable summary.
	Message string
}

// NewEvent constructs an Event for the given job and kind.
func NewEvent(kind EventKind, j *job.Job) Event {
	now := time.Now().UTC()
	return Event{
		Kind:       kind,
		Job:        j,
		ObservedAt: now,
		Message:    formatMessage(kind, j, now),
	}
}

func formatMessage(kind EventKind, j *job.Job, at time.Time) string {
	switch kind {
	case EventMissed:
		return fmt.Sprintf("job %q missed its scheduled run (detected at %s)",
			j.Name(), at.Format(time.RFC3339))
	case EventFailure:
		return fmt.Sprintf("job %q reported a failure (detected at %s)",
			j.Name(), at.Format(time.RFC3339))
	case EventRecovered:
		return fmt.Sprintf("job %q has recovered (detected at %s)",
			j.Name(), at.Format(time.RFC3339))
	default:
		return fmt.Sprintf("job %q: unknown event %q at %s",
			j.Name(), kind, at.Format(time.RFC3339))
	}
}
