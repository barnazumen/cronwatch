package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// JobConfig holds configuration for a single monitored cron job.
type JobConfig struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
}

// AlertConfig holds alerting backend settings.
type AlertConfig struct {
	Email   *EmailConfig   `yaml:"email,omitempty"`
	Webhook *WebhookConfig `yaml:"webhook,omitempty"`
}

// EmailConfig holds SMTP settings for email alerts.
type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
}

// WebhookConfig holds settings for webhook alerts.
type WebhookConfig struct {
	URL string `yaml:"url"`
}

// Config is the top-level application configuration.
type Config struct {
	LogLevel string      `yaml:"log_level"`
	Jobs     []JobConfig `yaml:"jobs"`
	Alerts   AlertConfig `yaml:"alerts"`
}

// Load reads and validates a YAML config file at the given path.
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
		return nil, fmt.Errorf("config: at least one job must be defined")
	}

	for i, j := range cfg.Jobs {
		if j.Name == "" {
			return nil, fmt.Errorf("config: job[%d] missing name", i)
		}
		if j.Schedule == "" {
			return nil, fmt.Errorf("config: job %q missing schedule", j.Name)
		}
	}

	return &cfg, nil
}
