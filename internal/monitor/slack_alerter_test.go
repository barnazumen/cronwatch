package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/job"
)

func makeSlackJob(t *testing.T) *job.Job {
	t.Helper()
	cfg := jobConfig("slack-job", "@hourly")
	j, err := job.NewJob(cfg)
	if err != nil {
		t.Fatalf("NewJob: %v", err)
	}
	j.RecordFailure("exit status 1")
	return j
}

func TestSlackAlerter_SendsOnFailure(t *testing.T) {
	var received slackPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	alerter := NewSlackAlerter(ts.URL, "#alerts")
	j := makeSlackJob(t)

	if err := alerter.Alert(j); err != nil {
		t.Fatalf("Alert returned error: %v", err)
	}

	if received.Channel != "#alerts" {
		t.Errorf("expected channel #alerts, got %q", received.Channel)
	}
	if received.Text == "" {
		t.Error("expected non-empty text in slack payload")
	}
}

func TestSlackAlerter_PropagatesError(t *testing.T) {
	alerter := NewSlackAlerter("http://127.0.0.1:0/no-server", "")
	j := makeSlackJob(t)

	if err := alerter.Alert(j); err == nil {
		t.Fatal("expected error from unreachable server, got nil")
	}
}

func TestSlackAlerter_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	alerter := NewSlackAlerter(ts.URL, "#alerts")
	j := makeSlackJob(t)

	if err := alerter.Alert(j); err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}
