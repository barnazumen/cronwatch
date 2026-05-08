package monitor_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/internal/monitor"
	"github.com/cronwatch/internal/watcher"
)

type countingAlerter struct {
	calls  int
	errs   []error
	events []*watcher.Event
}

func (c *countingAlerter) Alert(e *watcher.Event) error {
	c.calls++
	c.events = append(c.events, e)
	if len(c.errs) == 0 {
		return nil
	}
	err := c.errs[0]
	c.errs = c.errs[1:]
	return err
}

func makeRetryEvent(name string) *watcher.Event {
	return watcher.NewEvent(watcher.EventMissed, &fakeJob{name: name, schedule: "@hourly"}, time.Now())
}

func TestRetryAlerter_SucceedsFirstAttempt(t *testing.T) {
	inner := &countingAlerter{}
	ra := monitor.NewRetryAlerter(inner, 3, 0)
	e := makeRetryEvent("job-ok")

	if err := ra.Alert(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 1 {
		t.Errorf("expected 1 call, got %d", inner.calls)
	}
}

func TestRetryAlerter_RetriesOnTransientError(t *testing.T) {
	sentinel := errors.New("transient")
	inner := &countingAlerter{errs: []error{sentinel, sentinel}}
	ra := monitor.NewRetryAlerter(inner, 3, 0)
	e := makeRetryEvent("job-retry")

	if err := ra.Alert(e); err != nil {
		t.Fatalf("expected eventual success, got: %v", err)
	}
	if inner.calls != 3 {
		t.Errorf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryAlerter_ReturnsErrorAfterMaxAttempts(t *testing.T) {
	sentinel := errors.New("permanent")
	inner := &countingAlerter{errs: []error{sentinel, sentinel, sentinel}}
	ra := monitor.NewRetryAlerter(inner, 3, 0)
	e := makeRetryEvent("job-fail")

	if err := ra.Alert(e); err == nil {
		t.Fatal("expected error after max retries")
	}
	if inner.calls != 3 {
		t.Errorf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryAlerter_DefaultsApplied(t *testing.T) {
	inner := &countingAlerter{}
	ra := monitor.NewRetryAlerter(inner, 0, 0)
	e := makeRetryEvent("job-defaults")

	if err := ra.Alert(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls < 1 {
		t.Errorf("expected at least 1 call, got %d", inner.calls)
	}
}
