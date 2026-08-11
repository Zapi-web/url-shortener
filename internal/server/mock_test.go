package server

import (
	"context"
	"time"
)

type mockShortener struct {
	createFunc func(ctx context.Context, longURL string, userID uint64) (string, error)
	getFunc    func(ctx context.Context, shortURL string) (string, error)
}

func (m *mockShortener) Create(ctx context.Context, longURL string, userID uint64) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, longURL, userID)
	}

	return "", nil
}

func (m *mockShortener) Get(ctx context.Context, shortURL string) (string, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, shortURL)
	}

	return "", nil
}

type mockMetrics struct {
	incTotalCalls []incTotalCall
	observeCalls  []observeCall
}

type incTotalCall struct {
	handler string
	status  int
}

type observeCall struct {
	handler  string
	status   int
	duration time.Duration
}

func (m *mockMetrics) IncTotalRequests(handler string, status int) {
	m.incTotalCalls = append(m.incTotalCalls, incTotalCall{
		handler: handler,
		status:  status,
	})
}

func (m *mockMetrics) ObserveRequestDuration(handler string, status int, duration time.Duration) {
	m.observeCalls = append(m.observeCalls, observeCall{
		handler:  handler,
		status:   status,
		duration: duration,
	})
}
