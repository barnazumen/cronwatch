package watcher

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestCircuit_DefaultsApplied(t *testing.T) {
	c := NewCircuit(0, 0)
	if c.maxFailures != 3 {
		t.Fatalf("want maxFailures=3, got %d", c.maxFailures)
	}
	if c.resetAfter != 5*time.Minute {
		t.Fatalf("want resetAfter=5m, got %v", c.resetAfter)
	}
}

func TestCircuit_InitiallyClosed(t *testing.T) {
	c := NewCircuit(3, time.Minute)
	if !c.Allow() {
		t.Fatal("new circuit should allow")
	}
	if c.State() != CircuitClosed {
		t.Fatalf("want Closed, got %v", c.State())
	}
}

func TestCircuit_OpensAfterMaxFailures(t *testing.T) {
	c := NewCircuit(3, time.Minute)
	c.RecordFailure()
	c.RecordFailure()
	if c.State() != CircuitClosed {
		t.Fatal("should still be closed after 2 failures")
	}
	c.RecordFailure()
	if c.State() != CircuitOpen {
		t.Fatalf("want Open after 3 failures, got %v", c.State())
	}
}

func TestCircuit_SuppressesWhenOpen(t *testing.T) {
	c := NewCircuit(1, time.Minute)
	c.RecordFailure()
	if c.Allow() {
		t.Fatal("open circuit should not allow")
	}
}

func TestCircuit_HalfOpenAfterReset(t *testing.T) {
	now := time.Now()
	c := NewCircuit(1, time.Minute)
	c.now = fixedNow(now)
	c.RecordFailure() // opens circuit

	// advance clock past reset window
	c.now = fixedNow(now.Add(2 * time.Minute))
	if !c.Allow() {
		t.Fatal("circuit should allow in half-open state")
	}
	if c.State() != CircuitHalfOpen {
		t.Fatalf("want HalfOpen, got %v", c.State())
	}
}

func TestCircuit_ClosesOnSuccessFromHalfOpen(t *testing.T) {
	now := time.Now()
	c := NewCircuit(1, time.Minute)
	c.now = fixedNow(now)
	c.RecordFailure()
	c.now = fixedNow(now.Add(2 * time.Minute))
	c.Allow() // transition to half-open
	c.RecordSuccess()
	if c.State() != CircuitClosed {
		t.Fatalf("want Closed after success, got %v", c.State())
	}
}

func TestCircuit_ReOpensOnFailureFromHalfOpen(t *testing.T) {
	now := time.Now()
	c := NewCircuit(1, time.Minute)
	c.now = fixedNow(now)
	c.RecordFailure()
	c.now = fixedNow(now.Add(2 * time.Minute))
	c.Allow() // transition to half-open
	c.RecordFailure()
	if c.State() != CircuitOpen {
		t.Fatalf("want Open after failure in half-open, got %v", c.State())
	}
}
