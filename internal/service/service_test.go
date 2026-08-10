package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
)

type testMocks struct {
	cache    *mockCache
	database *mockDatabase
	encoder  *mockEncoder
	kgs      *mockKGS
	metrics  *mockMetrics
}

func defaultMocks() testMocks {
	return testMocks{
		cache:    &mockCache{},
		database: &mockDatabase{},
		encoder:  &mockEncoder{},
		kgs:      &mockKGS{},
		metrics:  &mockMetrics{},
	}
}

const (
	exampleURL        = "https://example.com"
	exampleUserID     = uint64(1)
	generatedID       = uint64(123456789)
	expectedShortCode = "8m0Kx"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(m *testMocks)
		inputLongURL  string
		inputUserID   uint64
		wantOutput    string
		wantError     bool
		expectedError error
	}{
		{
			name:         "Normal request",
			inputLongURL: exampleURL,
			inputUserID:  exampleUserID,
			setup: func(m *testMocks) {
				m.kgs.generateFunc = func() uint64 {
					return generatedID
				}
				m.encoder.encodeFunc = func(id uint64) string {
					if id != generatedID {
						t.Errorf("encoder received wrong ID = %d; want %d", id, generatedID)
					}
					return expectedShortCode
				}
				m.database.setFunc = func(ctx context.Context, url *domain.URL) error {
					if url.ID != generatedID {
						t.Errorf("db.Set ID = %d; want %d", url.ID, generatedID)
					}
					if url.LongURL != exampleURL {
						t.Errorf("db.Set LongURL = %q; want %q", url.LongURL, exampleURL)
					}
					if url.UserID != exampleUserID {
						t.Errorf("db.Set UserID = %d; want %d", url.UserID, exampleUserID)
					}
					return nil
				}
				m.cache.setFunc = func(ctx context.Context, shortURL, longURL string) error {
					if shortURL != expectedShortCode {
						t.Errorf("cache.Set shortURL = %q; want %q", shortURL, expectedShortCode)
					}
					if longURL != exampleURL {
						t.Errorf("cache.Set LongURL = %q; want %q", longURL, exampleURL)
					}
					return nil
				}
			},
			wantOutput: expectedShortCode,
		},
		{
			name:         "Database error",
			inputLongURL: exampleURL,
			inputUserID:  exampleUserID,
			setup: func(m *testMocks) {
				m.database.setFunc = func(ctx context.Context, url *domain.URL) error {
					return errors.New("db set failure")
				}
			},
			wantError: true,
		},
		{
			name:          "Empty Input",
			inputLongURL:  "",
			inputUserID:   exampleUserID,
			wantError:     true,
			expectedError: domain.ErrInputisEmpty,
		},
		{
			name:          "Invalid UserID",
			inputLongURL:  exampleURL,
			inputUserID:   0,
			wantError:     true,
			expectedError: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := defaultMocks()

			if tt.setup != nil {
				tt.setup(&mocks)
			}

			testService := New(mocks.database, mocks.cache, mocks.kgs, mocks.encoder, mocks.metrics, 10*time.Second)

			output, err := testService.Create(t.Context(), tt.inputLongURL, tt.inputUserID)

			if tt.wantError {
				if err == nil {
					t.Fatalf("Create(%q, %d) expected error, got nil", tt.inputLongURL, tt.inputUserID)
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("Create(%q, %d) error = %v; want %v", tt.inputLongURL, tt.inputUserID, err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create(%q, %d) unexpected error = %v", tt.inputLongURL, tt.inputUserID, err)
			}

			if output != tt.wantOutput {
				t.Errorf("Create(%q, %d) = %q; want %q", tt.inputLongURL, tt.inputUserID, output, tt.wantOutput)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(m *testMocks)
		inputShortURL string
		wantOutput    string
		wantError     bool
		expectedError error
	}{
		{
			name: "Normal Request (cache miss)",
			setup: func(m *testMocks) {
				m.cache.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					if shortURL != expectedShortCode {
						t.Errorf("cache.Get ShortURL = %q; want %q", shortURL, expectedShortCode)
					}

					return "", domain.ErrUrlNotFound
				}
				m.encoder.decodeFunc = func(s string) (uint64, error) {
					if s != expectedShortCode {
						t.Errorf("Encoder.Decode EncodedID = %q; want %q", s, expectedShortCode)
					}
					return generatedID, nil
				}
				m.database.getFunc = func(ctx context.Context, id uint64) (domain.URL, error) {
					if id != generatedID {
						t.Errorf("Database.Get id = %d; want %d", id, generatedID)
					}
					return domain.URL{LongURL: exampleURL}, nil
				}
				m.cache.setFunc = func(ctx context.Context, shortURL, longURL string) error {
					if shortURL != expectedShortCode {
						t.Errorf("cache.Set shortURL = %q; want %q", shortURL, expectedShortCode)
					}
					if longURL != exampleURL {
						t.Errorf("cache.Set LongURL = %q; want %q", longURL, exampleURL)
					}
					return nil
				}
			},
			inputShortURL: expectedShortCode,
			wantOutput:    exampleURL,
		},
		{
			name: "Normal Request (cache hit)",
			setup: func(m *testMocks) {
				m.cache.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					if shortURL != expectedShortCode {
						t.Errorf("cache.Get ShortURL = %q; want %q", shortURL, expectedShortCode)
					}

					return exampleURL, nil
				}
				m.database.getFunc = func(ctx context.Context, id uint64) (domain.URL, error) {
					t.Fatalf("unexpected database.Get() use")
					return domain.URL{}, nil
				}
			},
			inputShortURL: expectedShortCode,
			wantOutput:    exampleURL,
		},
		{
			name: "Failed Request (not found)",
			setup: func(m *testMocks) {
				m.cache.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					return "", domain.ErrUrlNotFound
				}
				m.encoder.decodeFunc = func(s string) (uint64, error) {
					return generatedID, nil
				}
				m.database.getFunc = func(ctx context.Context, id uint64) (domain.URL, error) {
					return domain.URL{}, domain.ErrUrlNotFound
				}
			},
			inputShortURL: expectedShortCode,
			wantError:     true,
			expectedError: domain.ErrUrlNotFound,
		},
		{
			name:          "Empty input",
			inputShortURL: "",
			wantError:     true,
			expectedError: domain.ErrInputisEmpty,
		},
		{
			name:          "Failed decoding",
			inputShortURL: expectedShortCode,
			setup: func(m *testMocks) {
				m.cache.getFunc = func(ctx context.Context, shortURL string) (string, error) {
					return "", domain.ErrUrlNotFound
				}
				m.encoder.decodeFunc = func(s string) (uint64, error) {
					return 0, errors.New("decoding failed")
				}
				m.database.getFunc = func(ctx context.Context, id uint64) (domain.URL, error) {
					t.Fatalf("unexpected database.Get() use")
					return domain.URL{}, nil
				}
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := defaultMocks()

			if tt.setup != nil {
				tt.setup(&mocks)
			}

			testService := New(mocks.database, mocks.cache, mocks.kgs, mocks.encoder, mocks.metrics, 10*time.Second)

			output, err := testService.Get(t.Context(), tt.inputShortURL)

			if tt.wantError {
				if err == nil {
					t.Fatalf("Get(%q) expected error, got nil", tt.inputShortURL)
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("Get(%q) error = %v; want %v", tt.inputShortURL, err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Fatalf("Get(%q) unexpected error = %v", tt.inputShortURL, err)
			}

			if output != tt.wantOutput {
				t.Errorf("Get(%q) = %q; want %q", tt.inputShortURL, output, tt.wantOutput)
			}
		})
	}
}
