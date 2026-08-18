package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb redis.UniversalClient
	ttl time.Duration
}

func NewRedis(ctx context.Context, addrs []string, masterName string, password string, ttl time.Duration) (*Redis, error) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            addrs,
		MasterName:       masterName,
		Password:         password,
		SentinelPassword: password,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping a db %w", err)
	}

	return &Redis{
		rdb: rdb,
		ttl: ttl,
	}, nil
}

func (r *Redis) Set(ctx context.Context, shortURL string, longURL string) error {
	if shortURL == "" || longURL == "" {
		return domain.ErrInputisEmpty
	}

	err := r.rdb.Set(ctx, shortURL, longURL, r.ttl).Err()

	if err != nil {
		return fmt.Errorf("failed to set a key-value in database: %w", err)
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, shortURL string) (string, error) {
	if shortURL == "" {
		return "", domain.ErrInputisEmpty
	}

	val, err := r.rdb.Get(ctx, shortURL).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrUrlNotFound
		}

		return "", fmt.Errorf("failed to get a value from a db: %w", err)
	}

	return val, nil
}

func (r *Redis) Close() error {
	return r.rdb.Close()
}
