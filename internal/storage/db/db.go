package db

import (
	"context"
	"fmt"
	"log/slog"

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
	err := d.rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		return fmt.Errorf("failed to set a key-value in database: %w", err)
	}

	slog.Debug("key saved", "key", key)

	return nil
}

func (d *Database) Get(ctx context.Context, key string) (string, error) {
	val, err := d.rdb.Get(ctx, key).Result()

	if err != nil {
		return "", fmt.Errorf("failed to get a value from a db: %w", err)
	}

	slog.Debug("value getted", "key", key, "value", val)

	return val, nil
}

func (d *Database) Close() {
	d.rdb.Close()
}
