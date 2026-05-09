package watcher

import (
	"testing"
	"time"
)

func TestHealthCheck_InitiallyUnknown(t *testing.T) {
	hc := NewHealthCheck()
	_, ok := hc.Status("backup")
	if ok {
		t.Fatal("expected no status for unseen job")
	}
}

func TestHealthCheck_RecordHealthy(t *testing.T) {
	hc := NewHealthCheck()
	hc.RecordHealthy("backup")
	s, ok := hc.Status("backup")
	if !ok {
		t.Fatal("expected status to exist")
	}
	if !s.Healthy {
		t.Errorf("expected Healthy=true, got false")
	}
	if s.JobName != "backup" {
		t.Errorf("expected JobName=backup, got %q", s.JobName)
	}
}

func TestHealthCheck_RecordUnhealthy(t *testing.T) {
	hc := NewHealthCheck()
	hc.RecordUnhealthy("backup", "missed run")
	s, ok := hc.Status("backup")
	if !ok {
		t.Fatal("expected status to exist")
	}
	if s.Healthy {
		t.Errorf("expected Healthy=false, got true")
	}
	if s.Message != "missed run" {
		t.Errorf("unexpected message: %q", s.Message)
	}
}

func TestHealthCheck_OverwritesPreviousState(t *testing.T) {
	hc := NewHealthCheck()
	hc.RecordUnhealthy("backup", "missed run")
	hc.RecordHealthy("backup")
	s, _ := hc.Status("backup")
	if !s.Healthy {
		t.Errorf("expected Healthy=true after recovery")
	}
}

func TestHealthCheck_LastCheckTimestamp(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	hc := NewHealthCheck()
	hc.nowFn = func() time.Time { return now }
	hc.RecordHealthy("db-dump")
	s, _ := hc.Status("db-dump")
	if !s.LastCheck.Equal(now) {
		t.Errorf("expected LastCheck=%v, got %v", now, s.LastCheck)
	}
}

func TestHealthCheck_Snapshot(t *testing.T) {
	hc := NewHealthCheck()
	hc.RecordHealthy("job-a")
	hc.RecordUnhealthy("job-b", "timeout")
	snap := hc.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}

func TestHealthCheck_SnapshotIsIndependent(t *testing.T) {
	hc := NewHealthCheck()
	hc.RecordHealthy("job-a")
	snap := hc.Snapshot()
	snap[0].Healthy = false
	s, _ := hc.Status("job-a")
	if !s.Healthy {
		t.Error("snapshot mutation affected internal state")
	}
}
