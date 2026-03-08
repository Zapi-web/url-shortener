package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Zapi-web/url-shortener/internal/logger"
	"github.com/Zapi-web/url-shortener/internal/storage/redis"
)

func main() {
	ctx := context.Background()
	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}

	logLvl := os.Getenv("LOG_LEVEL")
	slog.SetDefault(logger.NewLogger(logLvl))
	slog.Info("Logger initialized", "level", logLvl)

	db := redis.NewDatabase(ctx, addr)
	slog.Info("Database initialized")

	err := db.Set(ctx, "Hello", "World")
	if err != nil {
		slog.Error("Failed to set a key and value in db", "err", err)
		return
	}

	val, err := db.Get(ctx, "Hello")
	if err != nil {
		slog.Error("Failed to get a value", "err", err)
	}

	slog.Info("test value", "val", val)
}
