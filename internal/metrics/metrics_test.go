package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestVictoriaMetrics_ExposeMetrics(t *testing.T) {
	vm := New()

	vm.IncTotalRequests("/api/v1/", http.StatusOK)
	vm.ObserveRequestDuration("/", http.StatusFound, 150*time.Millisecond)
	vm.IncTotalCacheRequest("hit")
	vm.ObserveQueryDuration("Get", 150*time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	vm.ExposeMetrics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ExposeMetrics() expected code = %d; got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()

	expectedStrings := []string{
		`requests_total{handler="/api/v1/",status="200"}`,
		`cache_requests_total{cache_status="hit"}`,
		`request_durations_seconds_bucket{handler="/",status="302"`,
		`query_duration_seconds_bucket{handler="Get"`,
	}

	for _, str := range expectedStrings {
		if !strings.Contains(body, str) {
			t.Errorf("ExposeMetrics() answer not contain: %v", str)
		}
	}
}
