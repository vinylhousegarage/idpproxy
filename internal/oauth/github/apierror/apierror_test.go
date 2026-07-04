package apierror

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")

	tests := []struct {
		name         string
		code         ErrorCode
		status       int
		err          error
		serverErrorArgs []string
		wantServerError []string
	}{
		{
			name:         "All arguments provided",
			code:         "TEST_CODE",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			serverErrorArgs: []string{"debug info"},
			wantServerError: []string{"debug info"},
		},
		{
			name:         "No server error info",
			code:         "TEST_CODE_NO_SERVER_ERROR",
			status:       http.StatusBadRequest,
			err:          originalErr,
			serverErrorArgs: nil,
			wantServerError: nil,
		},
		{
			name:         "Nil error",
			code:         "TEST_CODE_NIL_ERR",
			status:       http.StatusUnauthorized,
			err:          nil,
			serverErrorArgs: nil,
			wantServerError: nil,
		},
		{
			name:         "Multiple server error infos (ignores after first)",
			code:         "TEST_CODE_MULTI_SERVER_ERROR",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			serverErrorArgs: []string{"first info", "second info"},
			wantServerError: []string{"first info", "second info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.code, tt.status, tt.err, tt.serverErrorArgs...)

			if got.Code != tt.code {
				t.Errorf("expected code %q, got %q", tt.code, got.Code)
			}
			if got.HTTPStatus != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, got.HTTPStatus)
			}
			if got.Err != tt.err {
				t.Errorf("expected err %v, got %v", tt.err, got.Err)
			}
			if !reflect.DeepEqual(got.ServerError, tt.wantServerError) {
				t.Errorf("expected server error %v, got %v", tt.wantServerError, got.ServerError)
			}
		})
	}
}
