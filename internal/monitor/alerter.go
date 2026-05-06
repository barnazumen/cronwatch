package monitor

import (
	"fmt"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/job"
)

// Alerter is the interface implemented by all alert backends.
type Alerter interface {
	Alert(j *job.Job) error
}

// BuildAlerters constructs the list of Alerter instances from the application
// configuration. Unknown alerter types are silently skipped so that new types
// can be added without breaking existing deployments.
func BuildAlerters(cfg config.Config) []Alerter {
	var alerters []Alerter

	for _, a := range cfg.Alerters {
		switch a.Type {
		case "log":
			alerters = append(alerters, NewLogAlerter())
		case "email":
			alerters = append(alerters, NewEmailAlerter(
				a.Email.SMTPHost,
				a.Email.SMTPPort,
				a.Email.From,
				a.Email.To,
			))
		case "webhook":
			alerters = append(alerters, NewWebhookAlerter(a.Webhook.URL))
		case "slack":
			alerters = append(alerters, NewSlackAlerter(a.Slack.WebhookURL))
		case "pagerduty":
			alerters = append(alerters, NewPagerDutyAlerter(a.PagerDuty.IntegrationKey))
		default:
			fmt.Printf("cronwatch: unknown alerter type %q — skipping\n", a.Type)
		}
	}

	if len(alerters) == 0 {
		alerters = append(alerters, NewLogAlerter())
	}

	return alerters
}
