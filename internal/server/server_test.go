package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
)

const (
	exampleURL        = "https://example.com"
	exampleUserID     = uint64(1)
	expectedShortCode = "8m0Kx"
)

func TestHandlers_ServeHealthz(t *testing.T) {
	handlers := NewHandlers(&mockShortener{})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handlers.ServeHealthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ServeHealthz() expected status %d; got %d", http.StatusOK, rec.Code)
	}
}

func TestHandlers_ServeSaveUrl(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(m *mockShortener)
		inputRequest   saveRequest
		expectedStatus int
	}{
		{
			name: "Normal Request",
			setup: func(m *mockShortener) {
				m.createFunc = func(ctx context.Context, longURL string, userID uint64) (string, error) {
					if longURL != exampleURL {
						t.Fatalf("shortener.Create LongURL = %v; want %v", longURL, exampleURL)
					}
					if userID != exampleUserID {
						t.Fatalf("shortener.Create userID = %d; want %d", userID, exampleUserID)
					}

					return expectedShortCode, nil
				}
				m.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					t.Fatalf("unexpected shortener.Get call")
					return "", nil
				}
			},
			inputRequest: saveRequest{
				LongURL: exampleURL,
				UserID:  exampleUserID,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty input",
			setup: func(m *mockShortener) {
				m.createFunc = func(ctx context.Context, longURL string, userID uint64) (string, error) {
					t.Fatalf("unexpected shortener.Create call")
					return "", nil
				}
			},
			inputRequest: saveRequest{
				LongURL: "",
				UserID:  exampleUserID,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service error",
			setup: func(m *mockShortener) {
				m.createFunc = func(ctx context.Context, longURL string, userID uint64) (string, error) {
					return "", errors.New("internal error")
				}
			},
			inputRequest: saveRequest{
				LongURL: exampleURL,
				UserID:  exampleUserID,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testShortener := &mockShortener{}

			if tt.setup != nil {
				tt.setup(testShortener)
			}

			handlers := NewHandlers(testShortener)

			jsonBytes, err := json.Marshal(tt.inputRequest)
			if err != nil {
				t.Fatalf("unexpected marshal error: %v", err)
			}

			req := httptest.NewRequest("POST", "/api/v1/", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handlers.ServeSaveUrl(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ServeSaveUrl() expected code = %d; got %d", rec.Code, tt.expectedStatus)
			}
		})
	}
}

func TestHandlers_ServeGetURL(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(m *mockShortener)
		inputShortURL  string
		expectedURL    string
		expectedStatus int
	}{
		{
			name: "Normal Request",
			setup: func(m *mockShortener) {
				m.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					if shortURL != expectedShortCode {
						t.Fatalf("shortener.Get shortURL = %v; want %v", shortURL, expectedShortCode)
					}
					return exampleURL, nil
				}
				m.createFunc = func(ctx context.Context, longURL string, userID uint64) (string, error) {
					t.Fatalf("unexpected shortener.Create call")
					return "", nil
				}
			},
			inputShortURL:  expectedShortCode,
			expectedURL:    exampleURL,
			expectedStatus: http.StatusFound,
		},
		{
			name: "Empty Input",
			setup: func(m *mockShortener) {
				m.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					t.Fatalf("unexpected shortener.Get call")
					return "", nil
				}
			},
			inputShortURL:  "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			setup: func(m *mockShortener) {
				m.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					return "", domain.ErrUrlNotFound
				}
			},
			inputShortURL:  expectedShortCode,
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Internal Error",
			setup: func(m *mockShortener) {
				m.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					return "", errors.New("internal error")
				}
			},
			inputShortURL:  expectedShortCode,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testShortener := &mockShortener{}

			if tt.setup != nil {
				tt.setup(testShortener)
			}

			handlers := NewHandlers(testShortener)

			req := httptest.NewRequest("GET", "/"+tt.inputShortURL, nil)
			req.SetPathValue("id", tt.inputShortURL)
			rec := httptest.NewRecorder()

			handlers.ServeGetURL(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ServeGetURL() expected code = %d; got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusFound {
				gotUrl := rec.Header().Get("Location")

				if gotUrl != tt.expectedURL {
					t.Errorf("ServeGetURL() expected URL = %v; got %v", tt.expectedURL, gotUrl)
				}
			}
		})
	}
}

func TestMiddleware_MetricsMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handlerStatus  int
		expectedStatus int
		method         string
	}{
		{
			name:           "Default Status Code",
			handlerStatus:  0,
			expectedStatus: http.StatusOK,
			method:         "GET",
		},
		{
			name:           "Custom Status Code",
			handlerStatus:  http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
			method:         "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockMetrics{}
			middleware := NewMiddleware(mock)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Millisecond)
				if tt.handlerStatus != 0 {
					w.WriteHeader(tt.handlerStatus)
				}
			})

			wrappedHandler := middleware.MetricsMiddleware(nextHandler)

			req := httptest.NewRequest(tt.method, "/test", nil)
			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("MetricsMiddleware().IncTotalRequests expected code = %d; got %d", tt.expectedStatus, rec.Code)
			}

			if len(mock.incTotalCalls) != 1 {
				t.Fatalf("MetricsMiddleware().IncTotalRequests expected 1 call to metrics; got %d", len(mock.incTotalCalls))
			}
			incCall := mock.incTotalCalls[0]
			if incCall.handler != tt.method {
				t.Fatalf("MetricsMiddleware().IncTotalRequests expected method = %v; got %v", tt.method, incCall.handler)
			}
			if incCall.status != tt.expectedStatus {
				t.Fatalf("MetricsMiddleware().IncTotalRequests expected status in metrics = %d; got %d", tt.expectedStatus, incCall.status)
			}

			if len(mock.observeCalls) != 1 {
				t.Fatalf("MetricsMiddleware().ObserveRequestDuration expected 1 call to metrics; got %d", len(mock.incTotalCalls))
			}
			obsCall := mock.observeCalls[0]
			if obsCall.handler != tt.method {
				t.Errorf("MetricsMiddleware().ObserveRequestDuration expected method = %v; got %v", tt.method, incCall.handler)
			}
			if obsCall.status != tt.expectedStatus {
				t.Errorf("MetricsMiddleware().ObserveRequestDuration expected status in metrics = %d; got %d", tt.expectedStatus, incCall.status)
			}
			if obsCall.duration <= 0 {
				t.Errorf("MetricsMiddleware().ObserveRequestDuration expected duration > 0; got %d", obsCall.duration)
			}
		})
	}
}

func TestServer_LifecycleAndShutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate free port: %v", err)
	}

	port := fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
	_ = listener.Close()

	handlers := NewHandlers(&mockShortener{})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlers.ServeHealthz)

	srv := NewServer(port, mux, 1*time.Second, 1*time.Second, 1*time.Second)

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunServer()
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/healthz", port))
	if err != nil {
		t.Fatalf("Server; failed to send request on running server: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Server; expected status 200, got %d", resp.StatusCode)
	}

	if err := srv.Shutdown(); err != nil {
		t.Fatalf("Server.Shutdown unexpected error: %v", err)
	}

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("RunServer returned unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunServer did not return after Shutdown")
	}
}
