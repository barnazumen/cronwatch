package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cronwatch/internal/job"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	JobName   string    `json:"job_name"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// WebhookAlerter sends alert notifications to an HTTP webhook.
type WebhookAlerter struct {
	URL    string
	client *http.Client
}

// NewWebhookAlerter creates a WebhookAlerter that POSTs to the given URL.
func NewWebhookAlerter(url string) *WebhookAlerter {
	return &WebhookAlerter{
		URL: url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Alert serialises the job state as JSON and POSTs it to the webhook URL.
func (w *WebhookAlerter) Alert(j *job.Job) error {
	payload := WebhookPayload{
		JobName:   j.Name,
		Status:    j.Status.String(),
		Message:   fmt.Sprintf("cronwatch: job %q is in state %s", j.Name, j.Status.String()),
		Timestamp: time.Now().UTC(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook alerter: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook alerter: POST %s: %w", w.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook alerter: unexpected status %d from %s", resp.StatusCode, w.URL)
	}

	return nil
}
