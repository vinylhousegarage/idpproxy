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

func TestAPIError_GetHTTPStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  *APIError
		want int
	}{
		{
			name: "Priority given to own HTTPStatus",
			err:  &APIError{HTTPStatus: 404},
			want: 404,
		},
		{
			name: "Fallback to first Internal error status when HTTPStatus is 0",
			err: &APIError{
				HTTPStatus: 0,
				Internal:   []Internal{{Status: 400}, {Status: 403}},
			},
			want: 400,
		},
		{
			name: "Default to 500 when neither HTTPStatus nor Internal is set",
			err:  &APIError{HTTPStatus: 0, Internal: nil},
			want: 500,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.err.GetHTTPStatus()
			if got != tt.want {
				t.Errorf("GetHTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
