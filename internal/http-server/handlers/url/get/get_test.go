package get_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/get"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type mockGetter struct {
	fn func(ctx context.Context, key string) (string, error)
}

func (m *mockGetter) Get(ctx context.Context, key string) (string, error) {
	return m.fn(ctx, key)
}

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name       string
		shortURL   string
		mockResp   string
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			shortURL:   "Good_URL",
			mockResp:   "https://google.com",
			mockErr:    nil,
			wantStatus: http.StatusFound,
		},
		{
			name:       "Not found",
			shortURL:   "unknown",
			mockErr:    domain.ErrUrlNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Internal server Error",
			shortURL:   "any_url",
			mockErr:    errors.New("unexpected db error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Empty Alias",
			mockErr:    domain.ErrKeyisEmpty,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock := &mockGetter{
				fn: func(ctx context.Context, key string) (string, error) {
					return tt.mockResp, tt.mockErr
				},
			}

			handler := get.GetNew(mock)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("short_url", tt.shortURL)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantStatus == http.StatusFound {
				require.Equal(t, tt.mockResp, rr.Header().Get("Location"))
			}
		})
	}
}
