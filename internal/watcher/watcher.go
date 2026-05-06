package watcher

import (
	"log"
	"time"

	"github.com/cronwatch/internal/job"
)

// Watcher periodically checks jobs for missed runs and updates their status.
type Watcher struct {
	registry *job.Registry
	interval time.Duration
	notify   chan<- *job.Job
	stop     chan struct{}
}

// New creates a new Watcher that checks the registry on the given interval.
// Missed jobs are sent to the notify channel.
func New(registry *job.Registry, interval time.Duration, notify chan<- *job.Job) *Watcher {
	return &Watcher{
		registry: registry,
		interval: interval,
		notify:   notify,
		stop:     make(chan struct{}),
	}
}

// Start begins the watch loop in a goroutine.
func (w *Watcher) Start() {
	go w.run()
}

// Stop signals the watcher to halt.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.check(time.Now())
		case <-w.stop:
			log.Println("watcher: stopped")
			return
		}
	}
}

func (w *Watcher) check(now time.Time) {
	for _, j := range w.registry.All() {
		if isMissed(j, now) {
			j.RecordMissed()
			w.notify <- j
		}
	}
}

// isMissed returns true when a job's deadline has passed and it has not run recently.
func isMissed(j *job.Job, now time.Time) bool {
	deadline := j.LastRun().Add(j.MaxInterval())
	return now.After(deadline)
}
