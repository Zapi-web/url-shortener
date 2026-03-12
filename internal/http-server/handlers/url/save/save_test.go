package save_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"github.com/Zapi-web/url-shortener/internal/http-server/handlers/url/save"
	"github.com/stretchr/testify/require"
)

type mockSetter struct {
	fn func(ctx context.Context, key, value string) error
}

func (m *mockSetter) Set(ctx context.Context, key, value string) error {
	return m.fn(ctx, key, value)
}

func TestSetHandler(t *testing.T) {
	tests := []struct {
		name       string
		inputURL   string
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			inputURL:   "https://google.com",
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid URL",
			inputURL:   "its-not-a-url",
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Empty URL",
			inputURL:   "",
			mockErr:    domain.ErrInputisEmpty,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Internal server Error",
			inputURL:   "https://google.com",
			mockErr:    errors.New("some db error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Conflict",
			inputURL:   "https://google.com",
			mockErr:    domain.ErrKeyAlreadyExist,
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSetter{
				fn: func(ctx context.Context, key, value string) error {
					return tt.mockErr
				},
			}

			handler := save.New(mock)

			input := fmt.Sprintf(`{"url":"%s"}`, tt.inputURL)

			req := httptest.NewRequest(http.MethodPost, "/save", strings.NewReader(input))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantStatus == http.StatusOK {
				var resp save.Response

				err := json.Unmarshal(rr.Body.Bytes(), &resp)

				require.NoError(t, err)
				require.NotEmpty(t, resp.ShortURL)
			}
		})
	}
}
