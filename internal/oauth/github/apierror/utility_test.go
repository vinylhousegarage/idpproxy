package apierror

import "testing"

func TestFormatDetail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		value   string
		want    string
		wantErr bool
	}{
		{"Success: normal input", "query_state", "12345", "query_state: 12345", false},
		{"Error: empty key", "", "value", "", true},
		{"Error: empty value", "key", "", "", true},
		{"Error: both empty", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := FormatDetail(tt.key, tt.value)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr && result != tt.want {
				t.Errorf("expected %s, got %s", tt.want, result)
			}
		})
	}
}
