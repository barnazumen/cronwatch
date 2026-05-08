package monitor

import (
	"fmt"
	"time"

	"github.com/cronwatch/internal/watcher"
)

// RetryAlerter wraps an Alerter and retries transient failures up to
// maxAttempts times, sleeping delay between each attempt.
type RetryAlerter struct {
	inner      Alerter
	maxAttempts int
	delay      time.Duration
}

// NewRetryAlerter returns a RetryAlerter that retries up to maxAttempts times
// with the given delay between attempts. If maxAttempts < 1 it defaults to 3;
// negative delay is treated as zero.
func NewRetryAlerter(inner Alerter, maxAttempts int, delay time.Duration) *RetryAlerter {
	if maxAttempts < 1 {
		maxAttempts = 3
	}
	if delay < 0 {
		delay = 0
	}
	return &RetryAlerter{
		inner:       inner,
		maxAttempts: maxAttempts,
		delay:       delay,
	}
}

// Alert calls the inner Alerter, retrying on error up to maxAttempts times.
// It returns the last error if all attempts fail.
func (r *RetryAlerter) Alert(e *watcher.Event) error {
	var lastErr error
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		if err := r.inner.Alert(e); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt < r.maxAttempts && r.delay > 0 {
			time.Sleep(r.delay)
		}
	}
	return fmt.Errorf("alert failed after %d attempts: %w", r.maxAttempts, lastErr)
}
