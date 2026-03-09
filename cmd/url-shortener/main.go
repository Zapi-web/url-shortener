package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Zapi-web/url-shortener/internal/logger"
	"github.com/Zapi-web/url-shortener/internal/storage/redis"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	ctx := context.Background()
	addr := os.Getenv("REDIS_ADDR")
	port := os.Getenv("PORT")

	if addr == "" {
		addr = "localhost:6379"
	}
	if port == "" {
		port = "8282"
	}

	logLvl := os.Getenv("LOG_LEVEL")
	slog.SetDefault(logger.NewLogger(logLvl))
	slog.Info("Logger initialized", "level", logLvl)

	db := redis.NewDatabase(ctx, addr)
	defer db.Close()

	slog.Info("Database initialized")

	r.Post("/save", save.New(db))

	slog.Info("Starting server", "addr", addr, "port", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "err", err)
	}
}
