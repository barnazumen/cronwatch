package watcher

import "time"

// Retry wraps an alert attempt with configurable retry logic.
// It will re-attempt the action up to MaxAttempts times, waiting
// WaitBetween between each attempt. If all attempts fail the last
// error is returned.
type Retry struct {
	MaxAttempts int
	WaitBetween time.Duration
	sleep       func(time.Duration)
}

// NewRetry returns a Retry with sensible defaults.
// MaxAttempts defaults to 3, WaitBetween defaults to 2 seconds.
func NewRetry(maxAttempts int, waitBetween time.Duration) *Retry {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	if waitBetween <= 0 {
		waitBetween = 2 * time.Second
	}
	return &Retry{
		MaxAttempts: maxAttempts,
		WaitBetween: waitBetween,
		sleep:       time.Sleep,
	}
}

// Do executes fn up to MaxAttempts times. It returns nil as soon as fn
// succeeds, or the last error if every attempt fails.
func (r *Retry) Do(fn func() error) error {
	var err error
	for i := 0; i < r.MaxAttempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < r.MaxAttempts-1 {
			r.sleep(r.WaitBetween)
		}
	}
	return err
}
