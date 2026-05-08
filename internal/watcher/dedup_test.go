package watcher

import (
	"testing"
	"time"
)

func TestDedup_FirstAlertAllowed(t *testing.T) {
	d := NewDedup(5 * time.Minute)
	if d.IsDuplicate("backup", "missed") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestDedup_SecondAlertSuppressed(t *testing.T) {
	d := NewDedup(5 * time.Minute)
	d.IsDuplicate("backup", "missed") // record it
	if !d.IsDuplicate("backup", "missed") {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestDedup_DifferentKeyAllowed(t *testing.T) {
	d := NewDedup(5 * time.Minute)
	d.IsDuplicate("backup", "missed")
	if d.IsDuplicate("backup", "failure") {
		t.Fatal("different key should not be suppressed")
	}
}

func TestDedup_DifferentJobAllowed(t *testing.T) {
	d := NewDedup(5 * time.Minute)
	d.IsDuplicate("backup", "missed")
	if d.IsDuplicate("cleanup", "missed") {
		t.Fatal("different job should not be suppressed")
	}
}

func TestDedup_AllowsAfterWindowExpires(t *testing.T) {
	d := NewDedup(10 * time.Millisecond)
	d.IsDuplicate("backup", "missed")
	time.Sleep(20 * time.Millisecond)
	if d.IsDuplicate("backup", "missed") {
		t.Fatal("expected alert after window expiry to be allowed")
	}
}

func TestDedup_ResetAllowsImmediately(t *testing.T) {
	d := NewDedup(5 * time.Minute)
	d.IsDuplicate("backup", "missed")
	d.Reset("backup")
	if d.IsDuplicate("backup", "missed") {
		t.Fatal("expected alert after Reset to be allowed")
	}
}

func TestDedup_ResetDoesNotAffectOtherJobs(t *testing.T) {
	d := NewDedup(5 * time.Minute)
	d.IsDuplicate("backup", "missed")
	d.IsDuplicate("cleanup", "missed")
	d.Reset("backup")
	if !d.IsDuplicate("cleanup", "missed") {
		t.Fatal("cleanup should still be suppressed after backup reset")
	}
}

func TestDedup_DefaultWindowApplied(t *testing.T) {
	d := NewDedup(0) // should default to 5 minutes
	if d.window != 5*time.Minute {
		t.Fatalf("expected default window 5m, got %v", d.window)
	}
}
