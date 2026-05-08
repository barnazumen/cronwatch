package monitor

import (
	"context"
	"fmt"

	"github.com/cronwatch/cronwatch/internal/job"
	"github.com/cronwatch/cronwatch/internal/watcher"
)

// circuitAlerter wraps an Alerter with a per-job circuit breaker so that
// repeated downstream failures do not cause an alert storm.
type circuitAlerter struct {
	inner    Alerter
	circuits map[string]*watcher.Circuit
	max      int
	reset    interface{ String() string }
	newCB    func() *watcher.Circuit
}

// NewCircuitAlerter wraps inner with a circuit breaker that opens after
// maxFailures consecutive alert delivery errors and resets after the
// configured duration encoded in cb.
func NewCircuitAlerter(inner Alerter, cb *watcher.Circuit) Alerter {
	// cb is used as a factory template; each job gets its own instance
	// with the same parameters via the closure below.
	return &circuitAlerter{
		inner:    inner,
		circuits: make(map[string]*watcher.Circuit),
		newCB:    func() *watcher.Circuit { return watcher.NewCircuit(cb.MaxFailures(), cb.ResetAfter()) },
	}
}

func (a *circuitAlerter) getCircuit(name string) *watcher.Circuit {
	if cb, ok := a.circuits[name]; ok {
		return cb
	}
	cb := a.newCB()
	a.circuits[name] = cb
	return cb
}

// Alert forwards to the inner alerter when the circuit is closed or half-open.
func (a *circuitAlerter) Alert(ctx context.Context, j *job.Job, ev watcher.Event) error {
	cb := a.getCircuit(j.Name())
	if !cb.Allow() {
		return fmt.Errorf("circuit open for job %q: alert suppressed", j.Name())
	}
	if err := a.inner.Alert(ctx, j, ev); err != nil {
		cb.RecordFailure()
		return err
	}
	cb.RecordSuccess()
	return nil
}
