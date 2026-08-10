package base62

import (
	"errors"
	"math"
	"testing"

	"github.com/Zapi-web/url-shortener/internal/domain"
)

func TestEncode(t *testing.T) {
	encoder := New()

	tests := []struct {
		name       string
		inputNum   uint64
		wantOutput string
	}{
		{
			name:       "Normal Test Number",
			inputNum:   123456789,
			wantOutput: "8m0Kx",
		},
		{
			name:       "Max uint64 input",
			inputNum:   math.MaxUint64,
			wantOutput: "lYGhA16ahyf",
		},
		{
			name:       "Min input 0",
			inputNum:   0,
			wantOutput: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := encoder.Encode(tt.inputNum)

			if tt.wantOutput != output {
				t.Errorf("Encode() wantOutput = %v got %v", tt.wantOutput, output)
			}
		})
	}
}

func TestDecoder(t *testing.T) {
	encoder := New()

	tests := []struct {
		name          string
		inputString   string
		wantOutput    uint64
		wantError     bool
		expectedError error
	}{
		{
			name:        "Normal Test String",
			inputString: "8m0Kx",
			wantOutput:  123456789,
		},
		{
			name:        "Minimal Input",
			inputString: "0",
			wantOutput:  0,
		},
		{
			name:        "Max Input",
			inputString: "lYGhA16ahyf",
			wantOutput:  math.MaxUint64,
		},
		{
			name:          "Empty Input",
			inputString:   "",
			wantError:     true,
			expectedError: domain.ErrInputisEmpty,
		},
		{
			name:          "Invalid Character Error",
			inputString:   "8m0Kx!",
			wantError:     true,
			expectedError: ErrInvalidCharacter,
		},
		{
			name:          "Overflow Error",
			inputString:   "zzzzzzzzzzz",
			wantError:     true,
			expectedError: ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num, err := encoder.Decode(tt.inputString)

			if tt.wantError {
				if err == nil {
					t.Fatalf("Decode(%q) expected error, got nil", tt.inputString)
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("Decode(%q) error = %v; want %v", tt.inputString, err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Fatalf("Decode(%q) unexpected error = %v", tt.inputString, err)
			}

			if num != tt.wantOutput {
				t.Errorf("Decode(%q) = %d; want %d", tt.inputString, num, tt.wantOutput)
			}
		})
	}
}
