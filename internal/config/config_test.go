package config

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		testConfig  Config
		expectError bool
	}{
		{
			name: "Clear Config",
			testConfig: Config{
				Server: ServerConfig{
					Port:   "8080",
					NodeID: 1,
				},
			},
		},
		{
			name: "Invalid Non-Numeric Port",
			testConfig: Config{
				Server: ServerConfig{
					Port:   "bad-port",
					NodeID: 1,
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Out-Of-Bounds Port",
			testConfig: Config{
				Server: ServerConfig{
					Port:   "65536",
					NodeID: 1,
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Port 0",
			testConfig: Config{
				Server: ServerConfig{
					Port:   "0",
					NodeID: 1,
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Negative NodeID",
			testConfig: Config{
				Server: ServerConfig{
					Port:   "8080",
					NodeID: -1,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testConfig.Validate()

			if (err != nil) != tt.expectError {
				t.Errorf("Validate() error = %v, ExpectError = %v", err, tt.expectError)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name         string
		inputLevel   string
		wantLogLevel string
	}{
		{
			name:         "Normal Config(info)",
			inputLevel:   "info",
			wantLogLevel: "info",
		},
		{
			name:         "Normal Config(warn)",
			inputLevel:   "warn",
			wantLogLevel: "warn",
		},
		{
			name:         "Bad Config: Unknown LogLevel",
			inputLevel:   "bad-input",
			wantLogLevel: "info",
		},
		{
			name:         "Empty LogLevel Fallback",
			inputLevel:   "",
			wantLogLevel: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testConfig := Config{App: AppConfig{LogLevel: tt.inputLevel}}
			testConfig.Normalize()

			if testConfig.App.LogLevel != tt.wantLogLevel {
				t.Errorf("LogLevel = %q, want %q", testConfig.App.LogLevel, tt.wantLogLevel)
			}
		})
	}
}
