package monitor

import (
	"errors"
	"testing"

	"github.com/user/cronwatch/internal/job"
)

// stubAlerter is a test double that records calls and optionally returns an error.
type stubAlerter struct {
	called bool
	err    error
}

func (s *stubAlerter) Alert(_ *job.Job) error {
	s.called = true
	return s.err
}

func makeJob(t *testing.T) *job.Job {
	t.Helper()
	j, err := job.NewJob("test-job", "@hourly")
	if err != nil {
		t.Fatalf("NewJob: %v", err)
	}
	return j
}

func TestMultiAlerter_CallsAll(t *testing.T) {
	a1 := &stubAlerter{}
	a2 := &stubAlerter{}
	ma := NewMultiAlerter(a1, a2)

	if err := ma.Alert(makeJob(t)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a1.called || !a2.called {
		t.Error("expected both alerters to be called")
	}
}

func TestMultiAlerter_ContinuesAfterError(t *testing.T) {
	sentinel := errors.New("alerter down")
	a1 := &stubAlerter{err: sentinel}
	a2 := &stubAlerter{}
	ma := NewMultiAlerter(a1, a2)

	err := ma.Alert(makeJob(t))
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel in error chain, got: %v", err)
	}
	if !a2.called {
		t.Error("second alerter should still be called after first fails")
	}
}

func TestMultiAlerter_CombinesErrors(t *testing.T) {
	e1 := errors.New("err1")
	e2 := errors.New("err2")
	ma := NewMultiAlerter(&stubAlerter{err: e1}, &stubAlerter{err: e2})

	err := ma.Alert(makeJob(t))
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !errors.Is(err, e1) || !errors.Is(err, e2) {
		t.Errorf("both errors should be present, got: %v", err)
	}
}

func TestMultiAlerter_Len(t *testing.T) {
	ma := NewMultiAlerter(&stubAlerter{}, &stubAlerter{}, &stubAlerter{})
	if ma.Len() != 3 {
		t.Errorf("expected Len 3, got %d", ma.Len())
	}
}

func TestMultiAlerter_Empty(t *testing.T) {
	ma := NewMultiAlerter()
	if err := ma.Alert(makeJob(t)); err != nil {
		t.Fatalf("empty MultiAlerter should not error: %v", err)
	}
}
