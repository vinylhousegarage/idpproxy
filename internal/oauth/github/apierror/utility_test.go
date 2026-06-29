package apierror

import (
	"testing"
)

func TestFormatDetail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{"Success: normal input", "query_state", "12345", "query_state: 12345"},
		{"Success: empty value", "key", "", "key: "},
		{"Error: empty key", "", "value", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDetail(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
