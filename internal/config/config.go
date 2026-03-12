package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Port     string `env:"PORT" env-default:"8282"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

func ConfigInit() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
