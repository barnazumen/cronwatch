package monitor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/job"
	"github.com/cronwatch/cronwatch/internal/watcher"
)

type countingAlerter struct {
	calls int
	err   error
}

func (c *countingAlerter) Alert(_ context.Context, _ *job.Job, _ watcher.Event) error {
	c.calls++
	return c.err
}

func makeCircuitJob(t *testing.T) *job.Job {
	t.Helper()
	j, err := job.NewJob(job.Config{Name: "circuit-job", Schedule: "@every 1m"})
	if err != nil {
		t.Fatalf("NewJob: %v", err)
	}
	return j
}

func TestCircuitAlerter_ForwardsWhenClosed(t *testing.T) {
	inner := &countingAlerter{}
	cb := watcher.NewCircuit(3, time.Minute)
	a := NewCircuitAlerter(inner, cb)
	j := makeCircuitJob(t)
	ev := watcher.Event{}

	if err := a.Alert(context.Background(), j, ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 1 {
		t.Fatalf("want 1 call, got %d", inner.calls)
	}
}

func TestCircuitAlerter_OpensAfterMaxFailures(t *testing.T) {
	sentinel := errors.New("downstream down")
	inner := &countingAlerter{err: sentinel}
	cb := watcher.NewCircuit(2, time.Minute)
	a := NewCircuitAlerter(inner, cb)
	j := makeCircuitJob(t)
	ev := watcher.Event{}
	ctx := context.Background()

	// two failures open the circuit
	a.Alert(ctx, j, ev) //nolint:errcheck
	a.Alert(ctx, j, ev) //nolint:errcheck

	// third call should be suppressed without hitting inner
	err := a.Alert(ctx, j, ev)
	if err == nil {
		t.Fatal("expected suppression error")
	}
	if inner.calls != 2 {
		t.Fatalf("inner should have been called exactly 2 times, got %d", inner.calls)
	}
}

func TestCircuitAlerter_ResetsOnSuccess(t *testing.T) {
	inner := &countingAlerter{}
	cb := watcher.NewCircuit(3, time.Minute)
	a := NewCircuitAlerter(inner, cb)
	j := makeCircuitJob(t)
	ev := watcher.Event{}
	ctx := context.Background()

	a.Alert(ctx, j, ev) //nolint:errcheck
	a.Alert(ctx, j, ev) //nolint:errcheck
	if err := a.Alert(ctx, j, ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 3 {
		t.Fatalf("want 3 calls, got %d", inner.calls)
	}
}
