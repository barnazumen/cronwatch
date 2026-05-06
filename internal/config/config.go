package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// JobConfig holds configuration for a single monitored cron job.
type JobConfig struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Grace    int    `yaml:"grace_minutes"`
}

// AlertConfig holds configuration for alert destinations.
type AlertConfig struct {
	Email   *EmailConfig   `yaml:"email,omitempty"`
	Webhook *WebhookConfig `yaml:"webhook,omitempty"`
	Slack   *SlackConfig   `yaml:"slack,omitempty"`
}

// EmailConfig holds SMTP alert settings.
type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
}

// WebhookConfig holds HTTP webhook alert settings.
type WebhookConfig struct {
	URL string `yaml:"url"`
}

// SlackConfig holds Slack incoming webhook alert settings.
type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
}

// Config is the top-level application configuration.
type Config struct {
	LogLevel string      `yaml:"log_level"`
	Jobs     []JobConfig `yaml:"jobs"`
	Alerts   AlertConfig `yaml:"alerts"`
}

// Load reads and validates a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if len(cfg.Jobs) == 0 {
		return nil, errors.New("config: no jobs defined")
	}

	for i, j := range cfg.Jobs {
		if j.Name == "" {
			return nil, fmt.Errorf("config: job[%d]: name is required", i)
		}
		if j.Schedule == "" {
			return nil, fmt.Errorf("config: job %q: schedule is required", j.Name)
		}
	}

	return &cfg, nil
}
