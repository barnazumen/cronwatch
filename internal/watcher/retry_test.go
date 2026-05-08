package watcher

import (
	"errors"
	"testing"
	"time"
)

func TestRetry_DefaultsApplied(t *testing.T) {
	r := NewRetry(0, 0)
	if r.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", r.MaxAttempts)
	}
	if r.WaitBetween != 2*time.Second {
		t.Errorf("expected WaitBetween=2s, got %v", r.WaitBetween)
	}
}

func TestRetry_SucceedsOnFirstAttempt(t *testing.T) {
	r := NewRetry(3, time.Millisecond)
	r.sleep = func(time.Duration) {}
	calls := 0
	err := r.Do(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetry_RetriesOnFailure(t *testing.T) {
	r := NewRetry(3, time.Millisecond)
	r.sleep = func(time.Duration) {}
	sentinel := errors.New("boom")
	calls := 0
	err := r.Do(func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetry_SucceedsOnSecondAttempt(t *testing.T) {
	r := NewRetry(3, time.Millisecond)
	r.sleep = func(time.Duration) {}
	sentinel := errors.New("transient")
	calls := 0
	err := r.Do(func() error {
		calls++
		if calls < 2 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestRetry_SleepsBetweenAttempts(t *testing.T) {
	r := NewRetry(3, 500*time.Millisecond)
	sleptFor := make([]time.Duration, 0)
	r.sleep = func(d time.Duration) { sleptFor = append(sleptFor, d) }
	_ = r.Do(func() error { return errors.New("fail") })
	if len(sleptFor) != 2 {
		t.Fatalf("expected 2 sleeps, got %d", len(sleptFor))
	}
	for _, d := range sleptFor {
		if d != 500*time.Millisecond {
			t.Errorf("expected 500ms sleep, got %v", d)
		}
	}
}
