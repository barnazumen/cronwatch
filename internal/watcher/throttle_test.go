package watcher

import (
	"testing"
	"time"
)

func TestThrottle_AllowsFirst(t *testing.T) {
	th := NewThrottle(time.Minute)
	if !th.Allow("job-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottle_SuppressesWithinInterval(t *testing.T) {
	now := time.Now()
	th := newThrottleWithClock(time.Minute, func() time.Time { return now })

	th.Allow("job-a")
	if th.Allow("job-a") {
		t.Fatal("expected second call within interval to be suppressed")
	}
}

func TestThrottle_AllowsAfterInterval(t *testing.T) {
	now := time.Now()
	th := newThrottleWithClock(time.Minute, func() time.Time { return now })

	th.Allow("job-a")
	now = now.Add(61 * time.Second)
	if !th.Allow("job-a") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestThrottle_IndependentKeys(t *testing.T) {
	th := NewThrottle(time.Minute)

	th.Allow("job-a")
	if !th.Allow("job-b") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestThrottle_ResetAllowsImmediately(t *testing.T) {
	now := time.Now()
	th := newThrottleWithClock(time.Minute, func() time.Time { return now })

	th.Allow("job-a")
	th.Reset("job-a")
	if !th.Allow("job-a") {
		t.Fatal("expected Allow after Reset to succeed")
	}
}

func TestThrottle_DefaultIntervalOnZero(t *testing.T) {
	th := NewThrottle(0)
	if th.interval != time.Minute {
		t.Fatalf("expected default interval of 1m, got %v", th.interval)
	}
}
