package watcher

import (
	"strings"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
	"github.com/cronwatch/cronwatch/internal/job"
)

func makeEventJob(t *testing.T, name string) *job.Job {
	t.Helper()
	cfg := config.JobConfig{
		Name:           name,
		Schedule:       "@every 1m",
		GracePeriodSec: 10,
	}
	j, err := job.NewJob(cfg)
	if err != nil {
		t.Fatalf("NewJob: %v", err)
	}
	return j
}

func TestNewEvent_Missed(t *testing.T) {
	j := makeEventJob(t, "backup")
	before := time.Now().UTC()
	ev := NewEvent(EventMissed, j)
	after := time.Now().UTC()

	if ev.Kind != EventMissed {
		t.Errorf("kind = %q, want %q", ev.Kind, EventMissed)
	}
	if ev.Job != j {
		t.Error("event job pointer mismatch")
	}
	if ev.ObservedAt.Before(before) || ev.ObservedAt.After(after) {
		t.Errorf("ObservedAt %v out of range [%v, %v]", ev.ObservedAt, before, after)
	}
	if !strings.Contains(ev.Message, "backup") {
		t.Errorf("message missing job name: %q", ev.Message)
	}
	if !strings.Contains(ev.Message, "missed") {
		t.Errorf("message missing 'missed': %q", ev.Message)
	}
}

func TestNewEvent_Failure(t *testing.T) {
	j := makeEventJob(t, "deploy")
	ev := NewEvent(EventFailure, j)

	if ev.Kind != EventFailure {
		t.Errorf("kind = %q, want %q", ev.Kind, EventFailure)
	}
	if !strings.Contains(ev.Message, "failure") {
		t.Errorf("message missing 'failure': %q", ev.Message)
	}
}

func TestNewEvent_Recovered(t *testing.T) {
	j := makeEventJob(t, "sync")
	ev := NewEvent(EventRecovered, j)

	if ev.Kind != EventRecovered {
		t.Errorf("kind = %q, want %q", ev.Kind, EventRecovered)
	}
	if !strings.Contains(ev.Message, "recovered") {
		t.Errorf("message missing 'recovered': %q", ev.Message)
	}
}

func TestNewEvent_UnknownKind(t *testing.T) {
	j := makeEventJob(t, "test-job")
	ev := NewEvent(EventKind("unknown"), j)

	if !strings.Contains(ev.Message, "unknown") {
		t.Errorf("message missing 'unknown': %q", ev.Message)
	}
}
