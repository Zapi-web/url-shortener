package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	rdb *redis.Client
}

func NewDatabase(ctx context.Context, addr string) (*Database, error) {
	var d Database

	d.rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	err := d.rdb.Ping(ctx).Err()

	if err != nil {
		return nil, fmt.Errorf("failed to ping a db %w", err)
	}

	return &d, nil
}

func (d *Database) Set(ctx context.Context, key, value string) error {
	if key == "" || value == "" {
		return domain.ErrInputisEmpty
	}

	err := d.rdb.SetArgs(ctx, key, value, redis.SetArgs{
		TTL:  240 * time.Hour,
		Mode: "NX",
	}).Err()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.ErrKeyAlreadyExist
		}
		return fmt.Errorf("failed to set a key-value in database: %w", err)
	}

	slog.Debug("key saved", "key", key)

	return nil
}

func (d *Database) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", domain.ErrInputisEmpty
	}

	val, err := d.rdb.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrUrlNotFound
		}

		return "", fmt.Errorf("failed to get a value from a db: %w", err)
	}

	slog.Debug("value fetched", "key", key, "value", val)

	return val, nil
}

func (d *Database) Close() {
	_ = d.rdb.Close()
}
