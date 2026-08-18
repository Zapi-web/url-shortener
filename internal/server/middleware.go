package server

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type Metrics interface {
	IncTotalRequests(handler string, status int)
	ObserveRequestDuration(handler string, status int, duration time.Duration)
}

type Middleware struct {
	Metrics Metrics
}

func NewMiddleware(metrics Metrics) *Middleware {
	return &Middleware{
		Metrics: metrics,
	}
}

func (m *Middleware) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		slog.DebugContext(r.Context(), "new request", "method", r.Method)

		start := time.Now()
		next.ServeHTTP(rw, r)
		reqDuration := time.Since(start)

		pattern := r.Pattern
		if pattern == "" {
			pattern = "unmatched"
		}

		m.Metrics.IncTotalRequests(pattern, rw.statusCode)
		m.Metrics.ObserveRequestDuration(pattern, rw.statusCode, reqDuration)
	})
}
