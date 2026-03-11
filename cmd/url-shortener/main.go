package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/get"
	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Zapi-web/url-shortener/internal/logger"
	"github.com/Zapi-web/url-shortener/internal/storage/db"
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

	db, err := db.NewDatabase(ctx, addr)
	if err != nil {
		slog.Error("Failed to create a database", "err", err)
		return
	}
	defer db.Close()

	slog.Info("Database initialized")

	r.Post("/save", save.New(db))
	r.Get("/{short_url:.{22}==}", get.GetNew(db))

	slog.Info("Starting server", "addr", addr, "port", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "err", err)
	}
}
