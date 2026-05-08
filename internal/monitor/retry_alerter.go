package monitor

import (
	"fmt"

	"github.com/cronwatch/cronwatch/internal/job"
	"github.com/cronwatch/cronwatch/internal/watcher"
)

// RetryAlerter wraps an Alerter and retries failed alert deliveries
// using the provided Retry policy.
type RetryAlerter struct {
	inner Alerter
	retry *watcher.Retry
}

// NewRetryAlerter returns a RetryAlerter that delegates to inner and
// retries delivery according to the given Retry policy.
func NewRetryAlerter(inner Alerter, r *watcher.Retry) *RetryAlerter {
	return &RetryAlerter{inner: inner, retry: r}
}

// Alert attempts to deliver the alert via the inner Alerter, retrying
// up to the configured number of times on failure.
func (ra *RetryAlerter) Alert(j *job.Job, event string) error {
	var lastErr error
	err := ra.retry.Do(func() error {
		lastErr = ra.inner.Alert(j, event)
		return lastErr
	})
	if err != nil {
		return fmt.Errorf("alert delivery failed after %d attempts: %w",
			ra.retry.MaxAttempts, err)
	}
	return nil
}
