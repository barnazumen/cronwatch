package job

import (
	"testing"
	"time"
)

func TestNewJob_Defaults(t *testing.T) {
	j := NewJob("backup", "0 2 * * *", 5*time.Minute)
	if j.Name != "backup" {
		t.Fatalf("expected name backup, got %s", j.Name)
	}
	if j.LastStatus != StatusUnknown {
		t.Fatalf("expected unknown status, got %s", j.LastStatus)
	}
	if j.FailCount != 0 {
		t.Fatalf("expected fail count 0, got %d", j.FailCount)
	}
}

func TestRecordSuccess(t *testing.T) {
	j := NewJob("sync", "* * * * *", time.Minute)
	now := time.Now()
	j.RecordSuccess(now)
	s := j.Snapshot()
	if s.LastStatus != StatusOK {
		t.Fatalf("expected ok, got %s", s.LastStatus)
	}
	if !s.LastSeen.Equal(now) {
		t.Fatalf("unexpected last seen time")
	}
	if s.FailCount != 0 {
		t.Fatalf("success should reset fail count")
	}
}

func TestRecordFailure(t *testing.T) {
	j := NewJob("report", "0 6 * * *", 10*time.Minute)
	now := time.Now()
	j.RecordFailure(now)
	j.RecordFailure(now)
	s := j.Snapshot()
	if s.LastStatus != StatusFailed {
		t.Fatalf("expected failed, got %s", s.LastStatus)
	}
	if s.FailCount != 2 {
		t.Fatalf("expected fail count 2, got %d", s.FailCount)
	}
}

func TestRecordMissed(t *testing.T) {
	j := NewJob("cleanup", "0 0 * * *", 15*time.Minute)
	j.RecordMissed()
	s := j.Snapshot()
	if s.LastStatus != StatusMissed {
		t.Fatalf("expected missed, got %s", s.LastStatus)
	}
	if s.FailCount != 1 {
		t.Fatalf("expected fail count 1, got %d", s.FailCount)
	}
}

func TestStatusString(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{StatusUnknown, "unknown"},
		{StatusOK, "ok"},
		{StatusFailed, "failed"},
		{StatusMissed, "missed"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}
