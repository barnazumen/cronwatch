package monitor

import (
	"fmt"
	"time"

	"github.com/cronwatch/cronwatch/internal/job"
	"github.com/cronwatch/cronwatch/internal/watcher"
)

// rateLimitedAlerter wraps an Alerter and suppresses duplicate alerts for the
// same job within a configurable cooldown window.
type rateLimitedAlerter struct {
	inner Alerter
	rl    *watcher.RateLimit
}

// NewRateLimitedAlerter wraps inner so that alerts for a given job are
// forwarded at most once per cooldown period.
func NewRateLimitedAlerter(inner Alerter, cooldown time.Duration) Alerter {
	return &rateLimitedAlerter{
		inner: inner,
		rl:    watcher.NewRateLimit(cooldown),
	}
}

func (r *rateLimitedAlerter) Alert(j *job.Job, reason string) error {
	if !r.rl.Allow(j.Name()) {
		return nil
	}
	if err := r.inner.Alert(j, reason); err != nil {
		// Roll back so the next tick retries.
		r.rl.Reset(j.Name())
		return fmt.Errorf("rate-limited alerter: %w", err)
	}
	return nil
}
