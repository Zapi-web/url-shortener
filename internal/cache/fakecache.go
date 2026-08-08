package cache

import (
	"context"

	"github.com/Zapi-web/url-shortener/internal/domain"
)

type FakeCache struct {
}

func NewFake() *FakeCache {
	return &FakeCache{}
}

func (f *FakeCache) Set(ctx context.Context, shortURL string, longURL string) error {
	return nil
}

func (f *FakeCache) Get(ctx context.Context, shortURL string) (string, error) {
	return "", domain.ErrUrlNotFound
}
