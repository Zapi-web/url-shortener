package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zapi-web/url-shortener/internal/config"
	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/get"
	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Zapi-web/url-shortener/internal/logger"
	"github.com/Zapi-web/url-shortener/internal/storage/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	ctx := context.Background()
	cfg, err := config.ConfigInit()

	if err != nil {
		slog.Error("Failed to read config", "err", err)
		return
	}

	slog.SetDefault(logger.NewLogger(cfg.LogLevel))
	slog.Info("Logger initialized", "level", cfg.LogLevel)

	db, err := db.NewDatabase(ctx, cfg.Addr)
	if err != nil {
		slog.Error("Failed to create a database", "err", err)
		return
	}
	defer db.Close()

	slog.Info("Database initialized")

	r.Post("/save", save.New(db))
	r.Get("/{short_url}", get.GetNew(db))

	slog.Info("Starting server", "addr", cfg.Addr, "port", cfg.Port)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverError:
		if err != nil {
			slog.Error("failed to start server", "err", err)
		}
	case sig := <-sign:
		slog.Info("Received a signal. Trying to gracefull shutdown", "sig", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			slog.Error("Could not stop server gracefully", "err", err)
		}
	}

	slog.Info("Server stopped")
}
