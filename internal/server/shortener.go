package server

import "context"

type Shortener interface {
	Create(ctx context.Context, longURL string, userID uint64) (string, error)
	Get(ctx context.Context, shortURL string) (string, error)
}
