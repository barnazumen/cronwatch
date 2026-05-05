package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// JobConfig holds configuration for a single monitored cron job.
type JobConfig struct {
	Name        string        `yaml:"name"`
	Schedule    string        `yaml:"schedule"`
	GracePeriod time.Duration `yaml:"grace_period"`
}

// Config is the top-level configuration structure.
type Config struct {
	LogLevel string      `yaml:"log_level"`
	Jobs     []JobConfig `yaml:"jobs"`
}

// Load reads and validates the configuration file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if len(cfg.Jobs) == 0 {
		return nil, fmt.Errorf("config must define at least one job")
	}

	for i, j := range cfg.Jobs {
		if j.Name == "" {
			return nil, fmt.Errorf("job[%d]: name is required", i)
		}
		if j.Schedule == "" {
			return nil, fmt.Errorf("job %q: schedule is required", j.Name)
		}
		if j.GracePeriod <= 0 {
			cfg.Jobs[i].GracePeriod = 5 * time.Minute
		}
	}

	return &cfg, nil
}
