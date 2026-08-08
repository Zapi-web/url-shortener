package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Server
	Port            string        `env:"PORT" env-default:"8080"`
	LogLevel        string        `env:"LOG_LEVEL" env-default:"info"`
	NodeID          int64         `env:"NODE_ID" env-required:"true"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"5s"`

	// Cache
	RedisAddr     string        `env:"REDIS_ADDR" env-default:"localhost:6379"`
	RedisPassword string        `env:"REDIS_PASSWORD" env-default:""`
	CacheTTL      time.Duration `env:"CACHE_TTL" env-default:"24h"`

	// Database
	DbTTL      time.Duration `env:"DATABASE_TTL" env-default:"336h"`
	ConnString string        `env:"CONNECTION_STRING" env-required:"true"`
}

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

	return &cfg, nil
}

func (c *Config) Validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid PORT %q: %w", c.Port, err)
	}

	if c.NodeID < 0 {
		return fmt.Errorf("invalid NODE_ID %d: %w", c.NodeID, err)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		slog.Warn("invalid LOG_LEVEL, fallback to 'info'", "LOG_LEVEL", c.LogLevel)
		c.LogLevel = "info"
	}

	return nil
}
