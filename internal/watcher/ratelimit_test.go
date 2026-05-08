package watcher

import (
	"testing"
	"time"
)

func TestRateLimit_AllowsFirstAlert(t *testing.T) {
	rl := NewRateLimit(5 * time.Minute)
	if !rl.Allow("backup") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestRateLimit_SuppressesWithinCooldown(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	rl := NewRateLimit(5 * time.Minute)
	rl.now = func() time.Time { return now }

	rl.Allow("backup") // first: allowed, records now

	rl.now = func() time.Time { return now.Add(2 * time.Minute) }
	if rl.Allow("backup") {
		t.Fatal("expected alert to be suppressed within cooldown")
	}
}

func TestRateLimit_AllowsAfterCooldown(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	rl := NewRateLimit(5 * time.Minute)
	rl.now = func() time.Time { return now }

	rl.Allow("backup")

	rl.now = func() time.Time { return now.Add(5 * time.Minute) }
	if !rl.Allow("backup") {
		t.Fatal("expected alert to be allowed after cooldown expires")
	}
}

func TestRateLimit_IndependentJobNames(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	rl := NewRateLimit(5 * time.Minute)
	rl.now = func() time.Time { return now }

	rl.Allow("backup")

	if !rl.Allow("cleanup") {
		t.Fatal("expected independent job to be allowed")
	}
}

func TestRateLimit_ResetAllowsImmediately(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	rl := NewRateLimit(5 * time.Minute)
	rl.now = func() time.Time { return now }

	rl.Allow("backup")
	rl.Reset("backup")

	rl.now = func() time.Time { return now.Add(1 * time.Second) }
	if !rl.Allow("backup") {
		t.Fatal("expected alert to be allowed after reset")
	}
}

func TestRateLimit_SnapshotIsIndependent(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	rl := NewRateLimit(5 * time.Minute)
	rl.now = func() time.Time { return now }
	rl.Allow("backup")

	snap := rl.Snapshot()
	delete(snap, "backup")

	if _, ok := rl.Snapshot()["backup"]; !ok {
		t.Fatal("modifying snapshot must not affect internal state")
	}
}
