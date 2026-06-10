package httpclient

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type dummyPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonBody    string
		wantError   bool
		expectedAge int
	}{
		{
			name:        "Success: valid JSON can be decoded",
			jsonBody:    `{"name": "Alice", "age": 30}`,
			wantError:   false,
			expectedAge: 30,
		},
		{
			name:        "Error: invalid JSON returns an error",
			jsonBody:    `{"name": "Bob", "age": }`,
			wantError:   true,
			expectedAge: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body := io.NopCloser(strings.NewReader(tt.jsonBody))
			resp := &http.Response{Body: body}

			var target dummyPayload
			err := DecodeJSON(resp, &target)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if target.Age != tt.expectedAge {
				t.Errorf("decoded result mismatch: got %v, want %v", target.Age, tt.expectedAge)
			}
		})
	}
}
