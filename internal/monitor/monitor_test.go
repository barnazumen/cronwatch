package monitor_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/job"
	"github.com/cronwatch/internal/monitor"
)

type fakeAlerter struct {
	mu     sync.Mutex
	alerts []alertCall
}

type alertCall struct {
	name    string
	status  job.Status
	message string
}

func (f *fakeAlerter) Alert(name string, status job.Status, message string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.alerts = append(f.alerts, alertCall{name, status, message})
	return nil
}

func (f *fakeAlerter) Calls() []alertCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]alertCall(nil), f.alerts...)
}

func makeRegistry(t *testing.T, jobs []config.JobConfig) *job.Registry {
	t.Helper()
	cfg := &config.Config{Jobs: jobs}
	reg, err := job.NewRegistry(cfg)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	return reg
}

func TestMonitor_AlertsOnMissedJob(t *testing.T) {
	cfgs := []config.JobConfig{
		{Name: "backup", Schedule: "@hourly", GracePeriod: "1s"},
	}
	reg := makeRegistry(t, cfgs)

	// Force the job to appear missed by setting last run far in the past.
	j, _ := reg.Get("backup")
	j.ForceLastRun(time.Now().Add(-2 * time.Hour))

	alerter := &fakeAlerter{}
	mon := monitor.New(reg, 50*time.Millisecond, alerter)
	mon.Start()
	time.Sleep(150 * time.Millisecond)
	mon.Stop()

	calls := alerter.Calls()
	if len(calls) == 0 {
		t.Fatal("expected at least one alert, got none")
	}
	if calls[0].name != "backup" {
		t.Errorf("expected alert for 'backup', got %q", calls[0].name)
	}
	if calls[0].status != job.StatusMissed {
		t.Errorf("expected StatusMissed, got %v", calls[0].status)
	}
}

func TestMonitor_NoAlertWhenJobRunsOnTime(t *testing.T) {
	cfgs := []config.JobConfig{
		{Name: "cleanup", Schedule: "@hourly", GracePeriod: "10m"},
	}
	reg := makeRegistry(t, cfgs)

	j, _ := reg.Get("cleanup")
	j.ForceLastRun(time.Now())

	alerter := &fakeAlerter{}
	mon := monitor.New(reg, 50*time.Millisecond, alerter)
	mon.Start()
	time.Sleep(120 * time.Millisecond)
	mon.Stop()

	if calls := alerter.Calls(); len(calls) != 0 {
		t.Errorf("expected no alerts, got %d", len(calls))
	}
}
