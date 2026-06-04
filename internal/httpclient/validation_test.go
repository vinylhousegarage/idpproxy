package httpclient

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
	}{
		{"Success 200", 200, "ok", false},
		{"Fail 500", 500, "error", true},
		{"Truncate long body", 400, strings.Repeat("a", 300), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			gotStatus, err := ValidateResponse(resp)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateResponse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotStatus != tt.statusCode {
				t.Errorf("ValidateResponse() status = %d, want %d", gotStatus, tt.statusCode)
			}

			if tt.wantErr {
				if !strings.Contains(err.Error(), "upstream error:") {
					t.Errorf("error message missing prefix")
				}
			}
		})
	}
}
