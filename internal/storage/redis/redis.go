package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	rdb *redis.Client
}

func NewDatabase(ctx context.Context, addr string) *Database {
	var d Database

	d.rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx = context.Background()

	return &d
}

func (d *Database) Set(ctx context.Context, key, value string) error {
	err := d.rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		return fmt.Errorf("failed to set a key-value in database: %w", err)
	}

	return nil
}

func (d *Database) Get(ctx context.Context, key string) (string, error) {
	val, err := d.rdb.Get(ctx, key).Result()

	if err != nil {
		return "", fmt.Errorf("failed to get a value from a db: %w", err)
	}

	return val, nil
}
