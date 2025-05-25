package hash

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		hashed      string
		alg         Algorithm
		expectedErr error
	}{
		{
			name:   "valid MD5 hash",
			input:  "test",
			hashed: "098f6bcd4621d373cade4e832627b4f6",
			alg:    AlgorithmMD5,
		},
		{
			name:   "valid SHA256 hash",
			input:  "test",
			hashed: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			alg:    AlgorithmSHA256,
		},
		{
			name:        "empty input",
			input:       "",
			hashed:      "hash",
			alg:         AlgorithmMD5,
			expectedErr: ErrEmptyInput,
		},
		{
			name:        "empty hash",
			input:       "test",
			hashed:      "",
			alg:         AlgorithmMD5,
			expectedErr: ErrEmptyHash,
		},
		{
			name:        "invalid algorithm",
			input:       "test",
			hashed:      "hash",
			alg:         Algorithm(99),
			expectedErr: ErrUnsupportedAlgorithm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := New(tt.input, tt.hashed, tt.alg)

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if h.Input() != tt.input {
				t.Errorf("expected input %q, got %q", tt.input, h.Input())
			}

			if h.Hashed() != tt.hashed {
				t.Errorf("expected hash %q, got %q", tt.hashed, h.Hashed())
			}

			if h.Algorithm() != tt.alg {
				t.Errorf("expected algorithm %v, got %v", tt.alg, h.Algorithm())
			}
		})
	}
}

func TestIsValidAlgorithm(t *testing.T) {
	tests := []struct {
		name     string
		alg      Algorithm
		expected bool
	}{
		{"MD5", AlgorithmMD5, true},
		{"SHA256", AlgorithmSHA256, true},
		{"invalid", Algorithm(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidAlgorithm(tt.alg)
			if result != tt.expected {
				t.Errorf("expected %v for algorithm %v, got %v", tt.expected, tt.alg, result)
			}
		})
	}
}
