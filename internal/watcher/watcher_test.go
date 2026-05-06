package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/job"
	"github.com/cronwatch/internal/watcher"
)

func makeRegistry(t *testing.T, maxInterval time.Duration) *job.Registry {
	t.Helper()
	cfg := &config.Config{
		Jobs: []config.JobConfig{
			{Name: "test-job", MaxInterval: maxInterval.String()},
		},
	}
	reg, err := job.NewRegistry(cfg)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	return reg
}

func TestWatcher_DetectsMissedJob(t *testing.T) {
	reg := makeRegistry(t, 1*time.Millisecond)
	notify := make(chan *job.Job, 1)
	w := watcher.New(reg, 5*time.Millisecond, notify)
	w.Start()
	defer w.Stop()

	select {
	case j := <-notify:
		if j.Name() != "test-job" {
			t.Errorf("expected test-job, got %s", j.Name())
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for missed job notification")
	}
}

func TestWatcher_NoAlertWhenJobRunsOnTime(t *testing.T) {
	reg := makeRegistry(t, 1*time.Hour)
	notify := make(chan *job.Job, 1)
	w := watcher.New(reg, 5*time.Millisecond, notify)
	w.Start()
	defer w.Stop()

	// Job has a 1h max interval so it should not be flagged immediately.
	select {
	case j := <-notify:
		t.Errorf("unexpected missed notification for job %s", j.Name())
	case <-time.After(50 * time.Millisecond):
		// expected: no alert
	}
}

func TestWatcher_StopsCleanly(t *testing.T) {
	reg := makeRegistry(t, 1*time.Hour)
	notify := make(chan *job.Job, 1)
	w := watcher.New(reg, 10*time.Millisecond, notify)
	w.Start()
	w.Stop()
	// Give the goroutine time to exit; no panic expected.
	time.Sleep(30 * time.Millisecond)
}
