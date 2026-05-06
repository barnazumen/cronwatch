package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/job"
	"github.com/user/cronwatch/internal/config"
)

func makePagerDutyJob() *job.Job {
	cfg := config.JobConfig{
		Name:     "nightly-backup",
		Schedule: "0 2 * * *",
		GracePeriod: 5,
	}
	j, _ := job.NewJob(cfg)
	j.RecordFailure("exit status 1")
	return j
}

func TestPagerDutyAlerter_SendsOnFailure(t *testing.T) {
	var received pagerDutyPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a := &pagerDutyAlerter{
		integrationKey: "test-key-123",
		client:         ts.Client(),
	}
	// Override the URL for testing by temporarily replacing the constant.
	// We use a real server so we patch via a custom client redirect.
	_ = a // prevent unused warning; test below uses overridden URL approach

	// Build alerter pointing at test server.
	a2 := &pagerDutyAlerter{
		integrationKey: "test-key-123",
		client: &http.Client{
			Transport: &urlRewriteTransport{base: ts.URL},
		},
	}

	if err := a2.Alert(makePagerDutyJob()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "test-key-123" {
		t.Errorf("routing key = %q, want %q", received.RoutingKey, "test-key-123")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Source != "nightly-backup" {
		t.Errorf("source = %q, want nightly-backup", received.Payload.Source)
	}
}

func TestPagerDutyAlerter_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a := &pagerDutyAlerter{
		integrationKey: "key",
		client: &http.Client{
			Transport: &urlRewriteTransport{base: ts.URL},
		},
	}
	if err := a.Alert(makePagerDutyJob()); err == nil {
		t.Error("expected error for non-2xx response")
	}
}

// urlRewriteTransport rewrites all requests to a fixed base URL for testing.
type urlRewriteTransport struct {
	base string
}

func (u *urlRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Scheme = "http"
	req2.URL.Host = u.base[len("http://"):]
	return http.DefaultTransport.RoundTrip(req2)
}
