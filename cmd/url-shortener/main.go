package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	defer rdb.Close()

	err := rdb.Set(ctx, "Hello", "World", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "Hello").Result()

	if err != nil {
		panic(err)
	}
	fmt.Println("Hello", val)
}
