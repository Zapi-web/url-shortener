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
	MetricsEnable   bool          `env:"ENABLE_METRICS" env-default:"false"`

	// Cache
	RedisAddrs      []string      `env:"REDIS_ADDRS" env-default:"localhost:6379" env-separator:","`
	RedisMasterName string        `env:"REDIS_MASTER_NAME" env-default:""`
	RedisPassword   string        `env:"REDIS_PASSWORD" env-default:""`
	CacheTTL        time.Duration `env:"CACHE_TTL" env-default:"24h"`

	// Database
	DbTTL       time.Duration `env:"DATABASE_TTL" env-default:"336h"`
	ConnString  string        `env:"CONNECTION_STRING" env-required:"true"`
	AutoMigrate bool          `env:"AUTO_MIGRATE" env-default:"false"`
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

	cfg.Normalize()

	return &cfg, nil
}

func (c *Config) Validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return fmt.Errorf("port is not a number %q: %w", c.Port, err)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("port is out-of-bounds %d", port)
	}

	if c.NodeID < 0 {
		return fmt.Errorf("invalid NODE_ID %d: %w", c.NodeID, err)
	}

	return nil
}

func (c *Config) Normalize() {
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		slog.Warn("invalid LOG_LEVEL, fallback to 'info'", "LOG_LEVEL", c.LogLevel)
		c.LogLevel = "info"
	}
}
