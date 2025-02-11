package utils

import (
	"testing"
)

func TestConvertToType(t *testing.T) {
	t.Parallel()

	type AliasType int64

	type testCase struct {
		input    string
		expected AliasType
		hasError bool
	}

	tests := []testCase{
		// Valid cases
		{"123", 123, false},
		{"0", 0, false},
		{"-999", -999, false},

		// Invalid cases
		{"abc", 0, true},
		{"", 0, true},
		{"123.45", 0, true},
		{" 123", 0, true},
		{"999999999999999999999999999999", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result, err := ConvertToType[AliasType](tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
