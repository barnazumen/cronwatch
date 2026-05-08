package watcher

import (
	"testing"
	"time"
)

func TestBackoff_DefaultsApplied(t *testing.T) {
	b := NewBackoff(0, 0)
	if b.Base != time.Minute {
		t.Fatalf("expected base=1m, got %v", b.Base)
	}
	if b.Cap < b.Base {
		t.Fatalf("cap must be >= base")
	}
}

func TestBackoff_CapEnforcedWhenCapLessThanBase(t *testing.T) {
	b := NewBackoff(10*time.Minute, 1*time.Minute)
	if b.Cap != b.Base {
		t.Fatalf("expected cap to be clamped to base, got cap=%v base=%v", b.Cap, b.Base)
	}
}

func TestBackoff_Doubles(t *testing.T) {
	b := NewBackoff(time.Minute, time.Hour)

	got := b.Next()
	if got != time.Minute {
		t.Fatalf("first Next: want 1m, got %v", got)
	}

	got = b.Next()
	if got != 2*time.Minute {
		t.Fatalf("second Next: want 2m, got %v", got)
	}

	got = b.Next()
	if got != 4*time.Minute {
		t.Fatalf("third Next: want 4m, got %v", got)
	}
}

func TestBackoff_CapsAtMax(t *testing.T) {
	b := NewBackoff(30*time.Minute, time.Hour)

	b.Next() // 30 m
	got := b.Next() // would be 60 m — exactly at cap
	if got != time.Hour {
		t.Fatalf("want 1h, got %v", got)
	}

	got = b.Next() // should stay at cap
	if got != time.Hour {
		t.Fatalf("want 1h (capped), got %v", got)
	}
}

func TestBackoff_Reset(t *testing.T) {
	b := NewBackoff(time.Minute, time.Hour)
	b.Next()
	b.Next()
	b.Reset()

	if b.Current() != time.Minute {
		t.Fatalf("after Reset: want 1m, got %v", b.Current())
	}
}

func TestBackoff_CurrentDoesNotAdvance(t *testing.T) {
	b := NewBackoff(time.Minute, time.Hour)

	for i := 0; i < 5; i++ {
		if b.Current() != time.Minute {
			t.Fatalf("Current() must not advance state (call %d)", i)
		}
	}
}
