package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zapi-web/url-shortener/internal/base62"
	"github.com/Zapi-web/url-shortener/internal/cache"
	"github.com/Zapi-web/url-shortener/internal/config"
	"github.com/Zapi-web/url-shortener/internal/database"
	"github.com/Zapi-web/url-shortener/internal/kgs"
	"github.com/Zapi-web/url-shortener/internal/logger"
	"github.com/Zapi-web/url-shortener/internal/metrics"
	"github.com/Zapi-web/url-shortener/internal/server"
	"github.com/Zapi-web/url-shortener/internal/service"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Init()

	if err != nil {
		slog.Error("Failed to read config", "err", err)
		return 1
	}

	slog.SetDefault(logger.NewLogger(cfg.LogLevel))
	slog.Info("Logger initialized", "level", cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgresDB, err := database.NewPostgres(ctx, cfg.ConnString)
	if err != nil {
		slog.Error("failed to connect to a database", "err", err)
		return 1
	}
	defer postgresDB.Close()

	slog.Info("Database initialized")

	var appCache service.Cache

	redisCache, err := cache.NewRedis(ctx, cfg.RedisAddr, cfg.CacheTTL)
	if err != nil {
		slog.Warn("Failed to connect to a cache, fallback to only database pattern", "err", err)
		appCache = cache.NewFake()
	} else {
		appCache = redisCache
		defer redisCache.Close()
	}

	slog.Info("Cache initialized")

	kgs, err := kgs.New(cfg.NodeID)
	if err != nil {
		slog.Error("failed to make an kgs", "err", err)
		return 1
	}

	encode := base62.New()
	Vmetrics := metrics.New()

	shortener := service.New(postgresDB, appCache, kgs, encode, Vmetrics, cfg.DbTTL)

	handlers := server.NewHandlers(shortener)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healtz", handlers.ServeHealthz)
	mux.HandleFunc("POST /api/v1/", handlers.ServeSaveUrl)
	mux.HandleFunc("GET /{id}", handlers.ServeGetURL)
	mux.HandleFunc("GET /metrics", Vmetrics.ExposeMetrics)

	mv := server.NewMiddleware(Vmetrics)
	handler := mv.MetricsMiddleware(mux)

	httpServer := server.NewServer(cfg.Port, handler)
	serverError := httpServer.RunServer(ctx)

	select {
	case err := <-serverError:
		if err != nil {
			slog.Error("received an error from http server", "err", err)
			stop()
			return 1
		}
	case <-ctx.Done():
		slog.Info("received a signal, starting graceful shutdown")
	}

	<-serverError // when chan is closed, means that server is fully stoped

	slog.Debug("Server stopped")
	return 0
}
