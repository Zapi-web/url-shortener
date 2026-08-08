package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Cache
	RedisAddr string        `env:"REDIS_ADDR" env-default:"localhost:6379"`
	CacheTTL  time.Duration `env:"CACHE_TTL" env-default:"24h"`

	// Database
	DbTTL      time.Duration `env:"DATABASE_TTL" env-default:"336h"`
	ConnString string        `env:"CONNECTION_STRING" env-required:"true"`

	// App
	Port     string `env:"PORT" env-default:"8080"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	NodeID   int64  `env:"NODE_ID" env-required:"true"`
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

	return nil
}
