package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/internal/job"
)

func makeWebhookJob(t *testing.T, name string) *job.Job {
	t.Helper()
	j, err := job.NewJob(name, "* * * * *")
	if err != nil {
		t.Fatalf("NewJob: %v", err)
	}
	j.RecordFailure()
	return j
}

func TestWebhookAlerter_SendsOnFailure(t *testing.T) {
	var received WebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	alerter := NewWebhookAlerter(server.URL)
	j := makeWebhookJob(t, "backup")

	if err := alerter.Alert(j); err != nil {
		t.Fatalf("Alert returned error: %v", err)
	}

	if received.JobName != "backup" {
		t.Errorf("job_name: got %q, want %q", received.JobName, "backup")
	}
	if received.Status != "failed" {
		t.Errorf("status: got %q, want %q", received.Status, "failed")
	}
	if received.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestWebhookAlerter_PropagatesHTTPError(t *testing.T) {
	alerter := NewWebhookAlerter("http://127.0.0.1:0/dead")
	j := makeWebhookJob(t, "cleanup")

	if err := alerter.Alert(j); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}

func TestWebhookAlerter_Non2xxReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	alerter := NewWebhookAlerter(server.URL)
	j := makeWebhookJob(t, "report")

	if err := alerter.Alert(j); err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}
