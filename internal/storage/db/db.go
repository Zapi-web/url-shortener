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
	res, err := d.rdb.SetNX(ctx, key, value, 240*time.Hour).Result()

	if res == false {
		return domain.ErrKeyAlreadyExist
	}

	if err != nil {
		return fmt.Errorf("failed to set a key-value in database: %w", err)
	}

	slog.Debug("key saved", "key", key)

	return nil
}

func (d *Database) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", domain.ErrKeyisEmpty
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
	d.rdb.Close()
}
