package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

type VictoriaMetrics struct{}

func New() *VictoriaMetrics {
	return &VictoriaMetrics{}
}

func (m *VictoriaMetrics) IncTotalRequests(handler string, status int) {
	name := fmt.Sprintf(`requests_total{handler="%s",status="%d"}`, handler, status)
	metrics.GetOrCreateCounter(name).Inc()
}

func (m *VictoriaMetrics) ObserveRequestDuration(handler string, status int, duration time.Duration) {
	name := fmt.Sprintf(`request_durations_seconds{handler="%s",status="%d"}`, handler, status)
	metrics.GetOrCreateHistogram(name).Update(duration.Seconds())
}

func (m *VictoriaMetrics) IncTotalCacheRequest(cacheStatus string) { // miss or hit
	name := fmt.Sprintf(`cache_requests_total{cache_status="%s"}`, cacheStatus)
	metrics.GetOrCreateCounter(name).Inc()
}

func (m *VictoriaMetrics) ObserveQueryDuration(handler string, duration time.Duration) {
	name := fmt.Sprintf(`query_duration_seconds{handler="%s"}`, handler)
	metrics.GetOrCreateHistogram(name).Update(duration.Seconds())
}

func (m *VictoriaMetrics) ExposeMetrics(w http.ResponseWriter, req *http.Request) {
	metrics.WritePrometheus(w, true)
}
