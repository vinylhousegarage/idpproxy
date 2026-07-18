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

func TestAPIError_AddInternal(t *testing.T) {
	t.Parallel()

	t.Run("success: internal error is added correctly", func(t *testing.T) {
		t.Parallel()

		apiErr := &APIError{}
		code := ErrorCode("TEST_CODE")
		value := "something went wrong"

		res1 := apiErr.AddInternal(code, value)
		if len(res1.Internals) != 1 {
			t.Fatalf("expected 1 internal error, but got %d", len(res1.Internals))
		}

		res2 := res1.AddInternal(code, value)
		if len(res2.Internals) != 2 {
			t.Fatalf("expected 2 internal errors, but got %d", len(res2.Internals))
		}

		if res1 != res2 {
			t.Error("expected returned APIError to be the same instance")
		}

		if res2.Internals[0].Code != code {
			t.Errorf("expected code: %s, got: %s", code, res2.Internals[0].Code)
		}
		if res2.Internals[0].Err.Error() != value {
			t.Errorf("expected error message: %s, got: %s", value, res2.Internals[0].Err.Error())
		}
	})

	t.Run("failure: empty values are ignored", func(t *testing.T) {
		t.Parallel()

		apiErr := &APIError{}
		res := apiErr.AddInternal("", "")

		if len(res.Internals) != 0 {
			t.Error("internal error should not be added when arguments are empty")
		}
	})
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
				Internals:  []APIInternal{{Status: 400}, {Status: 403}},
			},
			want: 400,
		},
		{
			name: "Default to 500 when neither HTTPStatus nor Internal is set",
			err:  &APIError{HTTPStatus: 0, Internals: nil},
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
