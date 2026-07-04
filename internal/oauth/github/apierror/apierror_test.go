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
	mockInternalErr := errors.New("internal debug info")

	tests := []struct {
		name         string
		code         ErrorCode
		status       int
		err          error
		internalArgs []Internal
		wantInternal []Internal
	}{
		{
			name:         "All arguments provided",
			code:         "TEST_CODE",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []Internal{{Code: "INTERNAL_DEBUG_CODE", Err: mockInternalErr}},
			wantInternal: []Internal{{Code: "INTERNAL_DEBUG_CODE", Err: mockInternalErr}},
		},
		{
			name:         "No internal info",
			code:         "TEST_CODE_NO_INTERNAL",
			status:       http.StatusBadRequest,
			err:          originalErr,
			internalArgs: nil,
			wantInternal: nil,
		},
		{
			name:         "Nil error",
			code:         "TEST_CODE_NIL_ERR",
			status:       http.StatusUnauthorized,
			err:          nil,
			internalArgs: nil,
			wantInternal: nil,
		},
		{
			name:         "Multiple internal infos (ignores after first)",
			code:         "TEST_CODE_MULTI_INTERNAL",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []Internal{
				{Code: "FIRST_CODE", Err: mockInternalErr},
				{Code: "SECOND_CODE", Err: mockInternalErr},
			},
			wantInternal: []Internal{
				{Code: "FIRST_CODE", Err: mockInternalErr},
				{Code: "SECOND_CODE", Err: mockInternalErr},
			},
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
			if !reflect.DeepEqual(got.Internal, tt.wantInternal) {
				t.Errorf("expected internal %v, got %v", tt.wantInternal, got.Internal)
			}
		})
	}
}
