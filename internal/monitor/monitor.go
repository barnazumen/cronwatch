package monitor

import (
	"log"
	"time"

	"github.com/cronwatch/internal/job"
)

// Monitor periodically checks jobs for missed runs and triggers alerts.
type Monitor struct {
	registry *job.Registry
	interval time.Duration
	alerter  Alerter
	stop     chan struct{}
}

// Alerter is the interface for sending alerts.
type Alerter interface {
	Alert(jobName string, status job.Status, message string) error
}

// New creates a new Monitor with the given registry, check interval, and alerter.
func New(registry *job.Registry, interval time.Duration, alerter Alerter) *Monitor {
	return &Monitor{
		registry: registry,
		interval: interval,
		alerter:  alerter,
		stop:     make(chan struct{}),
	}
}

// Start begins the monitoring loop in a background goroutine.
func (m *Monitor) Start() {
	go m.run()
}

// Stop signals the monitoring loop to exit.
func (m *Monitor) Stop() {
	close(m.stop)
}

func (m *Monitor) run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.check()
		case <-m.stop:
			log.Println("monitor: stopping")
			return
		}
	}
}

func (m *Monitor) check() {
	for _, j := range m.registry.All() {
		if j.IsMissed() {
			j.RecordMissed()
			msg := "job has not run within its expected schedule"
			if err := m.alerter.Alert(j.Name(), job.StatusMissed, msg); err != nil {
				log.Printf("monitor: alert failed for job %s: %v", j.Name(), err)
			}
		}
	}
}
