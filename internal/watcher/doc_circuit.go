// Package watcher provides job monitoring primitives.
//
// # Circuit Breaker
//
// Circuit implements a three-state circuit breaker (Closed → Open → Half-Open)
// that wraps alert delivery to avoid hammering a downstream alerting service
// when it is unhealthy.
//
// States:
//
//   - Closed  – normal operation; every alert attempt is allowed.
//   - Open    – alerting is suppressed after maxFailures consecutive delivery
//     failures; the circuit stays open for resetAfter duration.
//   - HalfOpen – one probe attempt is allowed after the reset window elapses;
//     a success closes the circuit, a failure re-opens it.
//
// Usage:
//
//	cb := watcher.NewCircuit(3, 2*time.Minute)
//	if cb.Allow() {
//	    if err := alerter.Alert(ctx, job, event); err != nil {
//	        cb.RecordFailure()
//	    } else {
//	        cb.RecordSuccess()
//	    }
//	}
package watcher
