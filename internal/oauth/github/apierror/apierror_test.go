package apierror

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")
	internalInfo := "debug info"

	tests := []struct {
		name         string
		code         ErrorCode
		status       int
		err          error
		internalArgs []string
		wantInternal string
	}{
		{
			name:         "All arguments provided",
			code:         "TEST_CODE",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []string{internalInfo},
			wantInternal: internalInfo,
		},
		{
			name:         "No internal info",
			code:         "TEST_CODE_NO_INTERNAL",
			status:       http.StatusBadRequest,
			err:          originalErr,
			internalArgs: nil,
			wantInternal: "",
		},
		{
			name:         "Nil error",
			code:         "TEST_CODE_NIL_ERR",
			status:       http.StatusUnauthorized,
			err:          nil,
			internalArgs: nil,
			wantInternal: "",
		},
		{
			name:         "Multiple internal infos (ignores after first)",
			code:         "TEST_CODE_MULTI_INTERNAL",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []string{"first info", "second info"},
			wantInternal: "first info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.code, tt.status, tt.err, tt.internalArgs...)

			if got.Code != tt.code {
				t.Errorf("expected code %q, got %q", tt.code, got.Code)
			}
			if got.HTTPStatus != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, got.HTTPStatus)
			}
			if got.Err != tt.err {
				t.Errorf("expected err %v, got %v", tt.err, got.Err)
			}
			if got.Internal != tt.wantInternal {
				t.Errorf("expected internal %q, got %q", tt.wantInternal, got.Internal)
			}
		})
	}
}
