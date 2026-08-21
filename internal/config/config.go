package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

func Init() (*Config, error) {
	var cfg Config

	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, fmt.Errorf("failed to read .env file: %w", err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("failed to read environment variables: %w", err)
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validate failed: %w", err)
	}

	cfg.Normalize()

	return &cfg, nil
}

func (c *Config) Validate() error {
	port, err := strconv.Atoi(c.Server.Port)
	if err != nil {
		return fmt.Errorf("port is not a number %q: %w", c.Server.Port, err)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("port is out-of-bounds %d", port)
	}

	if c.Server.NodeID < 0 {
		return fmt.Errorf("invalid NODE_ID %d: %w", c.Server.NodeID, err)
	}

	return nil
}

func (c *Config) Normalize() {
	validLogLevels := map[string]struct{}{"debug": {}, "info": {}, "warn": {}, "error": {}}
	if _, ok := validLogLevels[c.App.LogLevel]; !ok {
		slog.Warn("invalid LOG_LEVEL, fallback to 'info'", "LOG_LEVEL", c.App.LogLevel)
		c.App.LogLevel = "info"
	}

	validEnvironments := map[string]struct{}{"test": {}, "dev": {}, "prod": {}}
	if _, ok := validEnvironments[c.App.Environment]; !ok {
		slog.Warn("invalid ENVIRONMENT, fallback to 'test", "ENVIRONMENT", c.App.Environment)
		c.App.Environment = "test"
	}
}
