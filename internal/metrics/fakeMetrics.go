package metrics

import (
	"time"
)

type FakeMetrics struct{}

func NewFake() *FakeMetrics {
	return &FakeMetrics{}
}

func (m *FakeMetrics) IncTotalRequests(handler string, status int)                               {}
func (m *FakeMetrics) ObserveRequestDuration(handler string, status int, duration time.Duration) {}
func (m *FakeMetrics) IncTotalCacheRequest(cacheStatus string)                                   {}
func (m *FakeMetrics) ObserveQueryDuration(handler string, duration time.Duration)               {}
