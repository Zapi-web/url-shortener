package config

import "time"

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	Tracer   TracerConfig
}

type AppConfig struct {
	AppName       string `env:"APP_NAME" env-default:"url-shortener"`
	Environment   string `env:"ENVIRONMENT" env-default:"test"`
	MetricsEnable bool   `env:"ENABLE_METRICS" env-default:"false"`
	LogLevel      string `env:"LOG_LEVEL" env-default:"info"`
}

type ServerConfig struct {
	Port            string        `env:"PORT" env-default:"8080"`
	NodeID          int64         `env:"NODE_ID" env-required:"true"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
}

type DatabaseConfig struct {
	DbTTL       time.Duration `env:"DATABASE_TTL" env-default:"336h"`
	ConnString  string        `env:"CONNECTION_STRING" env-required:"true"`
	AutoMigrate bool          `env:"AUTO_MIGRATE" env-default:"false"`
}

type CacheConfig struct {
	RedisAddrs      []string      `env:"REDIS_ADDRS" env-default:"localhost:6379" env-separator:","`
	RedisMasterName string        `env:"REDIS_MASTER_NAME" env-default:""`
	RedisPassword   string        `env:"REDIS_PASSWORD" env-default:""`
	CacheTTL        time.Duration `env:"CACHE_TTL" env-default:"24h"`
}

type TracerConfig struct {
	TracesEnable  bool    `env:"ENABLE_TRACES" env-default:"false"`
	CollectorAddr string  `env:"COLLECTOR_ADDR" env-default:""`
	Ratio         float64 `env:"TRACER_RATIO" env-default:"1"`
}
