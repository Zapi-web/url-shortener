package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zapi-web/url-shortener/internal/base62"
	"github.com/Zapi-web/url-shortener/internal/cache"
	"github.com/Zapi-web/url-shortener/internal/config"
	"github.com/Zapi-web/url-shortener/internal/database"
	"github.com/Zapi-web/url-shortener/internal/kgs"
	"github.com/Zapi-web/url-shortener/internal/logger"
	"github.com/Zapi-web/url-shortener/internal/metrics"
	"github.com/Zapi-web/url-shortener/internal/server"
	"github.com/Zapi-web/url-shortener/internal/service"
	"github.com/Zapi-web/url-shortener/internal/tracer"
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

	if cfg.AutoMigrate {
		if err = postgresDB.Migrate(ctx); err != nil {
			slog.Error("failed to perform migration", "err", err)
			return 1
		}
	}

	slog.Info("Database initialized")

	var appCache service.Cache

	redisCache, err := cache.NewRedis(ctx, cfg.RedisAddrs, cfg.RedisMasterName, cfg.RedisPassword, cfg.CacheTTL)
	if err != nil {
		slog.Warn("Failed to connect to a cache, fallback to only database pattern", "err", err)
		appCache = cache.NewFake()
	} else {
		appCache = redisCache
		defer func() {
			if err := redisCache.Close(); err != nil {
				slog.Error("failed to close redis", "err", err)
			}
		}()
	}

	slog.Info("Cache initialized")

	kgs, err := kgs.New(cfg.NodeID)
	if err != nil {
		slog.Error("failed to make an kgs", "err", err)
		return 1
	}

	encode := base62.New()

	var prodMetrics interface {
		service.Metrics
		server.Metrics
	}

	if cfg.MetricsEnable {
		prodMetrics = metrics.NewVM()
	} else {
		prodMetrics = metrics.NewFake()
	}

	if cfg.TracesEnable && cfg.CollectorAddr != "" && cfg.Ratio > 0 {
		shutdown, err := tracer.New(ctx, cfg.AppName, cfg.Environment, cfg.CollectorAddr, cfg.Ratio)

		if err != nil {
			slog.Error("failed to inialize tracer, failover to work without tracer", "err", err)
		}
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdown(shutdownCtx); err != nil {
				slog.Error("failed to shutdown otel tracer")
			}
		}()
	}

	shortener := service.New(postgresDB, appCache, kgs, encode, prodMetrics, cfg.DbTTL)
	handler := server.NewHandlers(shortener)
	mv := server.NewMiddleware(prodMetrics)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.ServeHealthz)
	mux.HandleFunc("POST /api/v1/", handler.ServeSaveUrl)
	mux.HandleFunc("GET /{id}", handler.ServeGetURL)
	if cfg.MetricsEnable {
		mux.HandleFunc("GET /metrics", metrics.ExposeMetrics)
	}

	var rootHandler http.Handler = mux
	if cfg.MetricsEnable {
		rootHandler = mv.MetricsMiddleware(mux)
	}

	httpServer := server.NewServer(cfg.Port, rootHandler, cfg.ReadTimeout, cfg.WriteTimeout, cfg.ShutdownTimeout)
	serverError := make(chan error, 1)
	go func() {
		serverError <- httpServer.RunServer()
	}()

	select {
	case err := <-serverError:
		if err != nil {
			slog.Error("received an error from http server", "err", err)
			stop()
			return 1
		}
	case <-ctx.Done():
		slog.Info("received a signal, starting graceful shutdown")
		if err := httpServer.Shutdown(); err != nil {
			slog.Error("failed to shutdown server", "err", err)
			return 1
		}
	}

	slog.Debug("Server stopped")
	return 0
}
