package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/cronwatch/internal/job"
)

// SlackAlerter sends alert notifications to a Slack webhook URL.
type SlackAlerter struct {
	webhookURL string
	channel    string
	client     *http.Client
}

type slackPayload struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text"`
}

// NewSlackAlerter creates a SlackAlerter that posts to the given Slack incoming webhook URL.
func NewSlackAlerter(webhookURL, channel string) *SlackAlerter {
	return &SlackAlerter{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Alert sends a Slack message describing the job alert.
func (s *SlackAlerter) Alert(j *job.Job) error {
	msg := FormatAlert(j)
	payload := slackPayload{
		Channel: s.channel,
		Text:    msg,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack alerter: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack alerter: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack alerter: unexpected status %d", resp.StatusCode)
	}

	return nil
}
