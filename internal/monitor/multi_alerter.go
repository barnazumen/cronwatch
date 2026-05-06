package monitor

import (
	"errors"
	"fmt"

	"github.com/user/cronwatch/internal/job"
)

// MultiAlerter fans out an alert to multiple Alerter implementations,
// collecting all errors and returning them as a single combined error.
type MultiAlerter struct {
	alerters []Alerter
}

// NewMultiAlerter constructs a MultiAlerter that dispatches to each of the
// provided alerters in order.
func NewMultiAlerter(alerters ...Alerter) *MultiAlerter {
	return &MultiAlerter{alerters: alerters}
}

// Alert sends the alert to every registered alerter. All alerters are
// attempted even if an earlier one returns an error. A combined error is
// returned when one or more alerters fail.
func (m *MultiAlerter) Alert(j *job.Job) error {
	var errs []error
	for _, a := range m.alerters {
		if err := a.Alert(j); err != nil {
			errs = append(errs, fmt.Errorf("%T: %w", a, err))
		}
	}
	return errors.Join(errs...)
}

// Len returns the number of alerters registered with this MultiAlerter.
func (m *MultiAlerter) Len() int {
	return len(m.alerters)
}
