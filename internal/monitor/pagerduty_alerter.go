package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/cronwatch/internal/job"
)

const pagerDutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

type pagerDutyAlerter struct {
	integrationKey string
	client         *http.Client
}

type pagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     pagerDutyDetail   `json:"payload"`
	Client      string            `json:"client"`
}

type pagerDutyDetail struct {
	Summary   string    `json:"summary"`
	Source    string    `json:"source"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
}

// NewPagerDutyAlerter creates an Alerter that sends events to PagerDuty.
func NewPagerDutyAlerter(integrationKey string) Alerter {
	return &pagerDutyAlerter{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *pagerDutyAlerter) Alert(j *job.Job) error {
	body := pagerDutyPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Client:      "cronwatch",
		Payload: pagerDutyDetail{
			Summary:   FormatAlert(j),
			Source:    j.Name,
			Severity:  "error",
			Timestamp: time.Now().UTC(),
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := p.client.Post(pagerDutyEventsURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: send event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
