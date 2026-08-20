package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	Database   Database
	Cache      Cache
	Encoder    Encoder
	KGS        KGS
	Metrics    Metrics
	tracer     trace.Tracer
	defaultTTL time.Duration
}

func New(db Database, cache Cache, kgs KGS, encoder Encoder, metrics Metrics, ttl time.Duration) *Service {
	return &Service{
		Database:   db,
		Cache:      cache,
		Encoder:    encoder,
		KGS:        kgs,
		Metrics:    metrics,
		tracer:     otel.Tracer("url-shortener/service"),
		defaultTTL: ttl,
	}
}

func (s *Service) Create(ctx context.Context, longURL string, userID uint64) (string, error) {
	ctx, span := s.tracer.Start(ctx, "Service.Create",
		trace.WithAttributes(attribute.Int64("user_id", int64(userID))),
	)
	defer span.End()

	if longURL == "" {
		return "", domain.ErrInputisEmpty
	}
	if userID <= 0 {
		return "", domain.ErrInvalidInput
	}

	id := s.KGS.Generate()

	url := domain.URL{
		ID:        id,
		UserID:    userID,
		LongURL:   longURL,
		ExpiredAt: time.Now().UTC().Add(s.defaultTTL),
	}

	now := time.Now()
	err := s.Database.Set(ctx, &url)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		slog.ErrorContext(ctx, "failed to save url to database", "id", id, "user_id", userID, "err", err)
		return "", fmt.Errorf("failed to set a url %d in database: %w", id, err)
	}

	s.Metrics.ObserveQueryDuration("create", time.Since(now))
	slog.DebugContext(ctx, "saved url to database", "id", url.ID, "user_id", userID)

	encodedID := s.Encoder.Encode(url.ID)

	cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)

	go func() {
		defer cancel()

		if err := s.Cache.Set(cacheCtx, encodedID, longURL); err != nil {
			slog.WarnContext(cacheCtx, "failed to set record in cache", "key", encodedID, "err", err)
			return
		}
		slog.DebugContext(cacheCtx, "set key-value in cache", "key", encodedID)
	}()

	return encodedID, nil
}

func (s *Service) Get(ctx context.Context, shortURL string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "Service.Get",
		trace.WithAttributes(attribute.String("short_url", shortURL)),
	)
	defer span.End()

	if shortURL == "" {
		return "", domain.ErrInputisEmpty
	}

	getCacheCtx, getCancel := context.WithTimeout(ctx, 1*time.Second)
	res, err := s.Cache.Get(getCacheCtx, shortURL)
	getCancel()

	if err == nil {
		slog.DebugContext(ctx, "retrieved url from cache", "key", shortURL)
		s.Metrics.IncTotalCacheRequest("hit")
		return res, nil
	}

	if errors.Is(err, domain.ErrUrlNotFound) {
		s.Metrics.IncTotalCacheRequest("miss")
	} else {
		slog.WarnContext(ctx, "cache lookup failed, falling back to database", "key", shortURL, "err", err)
		s.Metrics.IncTotalCacheRequest("error")
	}

	decodedID, err := s.Encoder.Decode(shortURL)

	if err != nil {
		return "", fmt.Errorf("failed to decode short url: %w", err)
	}

	start := time.Now()
	url, err := s.Database.Get(ctx, decodedID)

	if err != nil {
		if errors.Is(err, domain.ErrUrlNotFound) {
			slog.DebugContext(ctx, "url not found in database", "short_url", shortURL, "decoded_id", decodedID)
			return "", domain.ErrUrlNotFound
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to query database for short url", "short_url", shortURL, "decoded_id", decodedID, "err", err)
		return "", fmt.Errorf("failed to find short url in database: %w", err)
	}

	s.Metrics.ObserveQueryDuration("get", time.Since(start))

	cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
	go func() {
		defer cancel()
		if err := s.Cache.Set(cacheCtx, shortURL, url.LongURL); err != nil {
			slog.WarnContext(cacheCtx, "failed to set a key-value record to a cache", "key", shortURL, "err", err)
			return
		}

		slog.DebugContext(cacheCtx, "populated cache on miss fallback", "key", shortURL)
	}()

	return url.LongURL, nil
}
