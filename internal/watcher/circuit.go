// Package watcher provides job monitoring primitives.
package watcher

import (
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	CircuitClosed   CircuitState = iota // normal operation
	CircuitOpen                         // alerting suppressed
	CircuitHalfOpen                     // testing if alerting should resume
)

// Circuit implements a simple circuit breaker for alert suppression.
// After maxFailures consecutive alert failures the circuit opens and
// suppresses further alerts until resetAfter has elapsed.
type Circuit struct {
	mu          sync.Mutex
	state       CircuitState
	failures    int
	maxFailures int
	resetAfter  time.Duration
	openedAt    time.Time
	now         func() time.Time
}

// NewCircuit creates a Circuit that opens after maxFailures consecutive
// failures and resets after resetAfter.
func NewCircuit(maxFailures int, resetAfter time.Duration) *Circuit {
	if maxFailures <= 0 {
		maxFailures = 3
	}
	if resetAfter <= 0 {
		resetAfter = 5 * time.Minute
	}
	return &Circuit{
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
		now:         time.Now,
	}
}

// Allow returns true when the circuit permits an alert attempt.
func (c *Circuit) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch c.state {
	case CircuitOpen:
		if c.now().Sub(c.openedAt) >= c.resetAfter {
			c.state = CircuitHalfOpen
			return true
		}
		return false
	default:
		return true
	}
}

// RecordSuccess resets failure count and closes the circuit.
func (c *Circuit) RecordSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures = 0
	c.state = CircuitClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (c *Circuit) RecordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures++
	if c.failures >= c.maxFailures {
		c.state = CircuitOpen
		c.openedAt = c.now()
	}
}

// State returns the current CircuitState.
func (c *Circuit) State() CircuitState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}
