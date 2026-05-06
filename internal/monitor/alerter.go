package monitor

import (
	"fmt"
	"log/slog"

	"github.com/user/cronwatch/internal/config"
)

// Alerter is the interface implemented by every alerting backend.
type Alerter interface {
	Alert(j interface{ GetStatus() string }) error
}

// BuildAlerters constructs a MultiAlerter from the application config,
// enabling every backend that has been configured. At least a LogAlerter
// is always included so there is a no-op default.
func BuildAlerters(cfg *config.Config, logger *slog.Logger) (*MultiAlerter, error) {
	var alerters []Alerter

	// Log alerter is always active.
	alerters = append(alerters, NewLogAlerter(logger))

	if cfg.Alerting.Email.Enabled {
		a, err := NewEmailAlerter(cfg.Alerting.Email)
		if err != nil {
			return nil, fmt.Errorf("email alerter: %w", err)
		}
		alerters = append(alerters, a)
	}

	if cfg.Alerting.Webhook.URL != "" {
		alerters = append(alerters, NewWebhookAlerter(cfg.Alerting.Webhook))
	}

	if cfg.Alerting.Slack.WebhookURL != "" {
		alerters = append(alerters, NewSlackAlerter(cfg.Alerting.Slack))
	}

	if cfg.Alerting.PagerDuty.RoutingKey != "" {
		alerters = append(alerters, NewPagerDutyAlerter(cfg.Alerting.PagerDuty))
	}

	return NewMultiAlerter(alerters...), nil
}
