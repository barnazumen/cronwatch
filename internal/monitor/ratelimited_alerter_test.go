package monitor

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/job"
)

type countingAlerter struct {
	calls atomic.Int32
	err   error
}

func (c *countingAlerter) Alert(_ *job.Job, _ string) error {
	c.calls.Add(1)
	return c.err
}

func makeRLJob(t *testing.T) *job.Job {
	t.Helper()
	j, err := job.NewJob(job.Config{Name: "rl-job", Schedule: "@hourly", Timeout: time.Minute})
	if err != nil {
		t.Fatalf("NewJob: %v", err)
	}
	return j
}

func TestRateLimitedAlerter_ForwardsFirst(t *testing.T) {
	inner := &countingAlerter{}
	a := NewRateLimitedAlerter(inner, 5*time.Minute)
	j := makeRLJob(t)

	if err := a.Alert(j, "missed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestRateLimitedAlerter_SuppressesSecond(t *testing.T) {
	inner := &countingAlerter{}
	a := NewRateLimitedAlerter(inner, time.Hour)
	j := makeRLJob(t)

	_ = a.Alert(j, "missed")
	_ = a.Alert(j, "missed")

	if inner.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestRateLimitedAlerter_RetriesAfterInnerError(t *testing.T) {
	inner := &countingAlerter{err: errors.New("smtp down")}
	a := NewRateLimitedAlerter(inner, time.Hour)
	j := makeRLJob(t)

	err := a.Alert(j, "missed")
	if err == nil {
		t.Fatal("expected error from inner alerter")
	}

	// After inner error the rate-limit should be reset so the next call retries.
	inner.err = nil
	if err2 := a.Alert(j, "missed"); err2 != nil {
		t.Fatalf("expected retry to succeed, got %v", err2)
	}
	if inner.calls.Load() != 2 {
		t.Fatalf("expected 2 inner calls, got %d", inner.calls.Load())
	}
}
