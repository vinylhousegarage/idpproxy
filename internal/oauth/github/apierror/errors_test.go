package apierror

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIErrors(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")
	internalInfo := "debug info"

	tests := []struct {
		name           string
		fn             func(error, ...string) *APIError
		expectedCode   ErrorCode
		expectedStatus int
	}{
		{
			name:           "Internal",
			fn:             Internal,
			expectedCode:   ErrorCodeInternal,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ProxyCodeIssue",
			fn:             ProxyCodeIssue,
			expectedCode:   ErrorCodeProxyCodeIssue,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.fn(originalErr, internalInfo)

			if got.Code != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, got.Code)
			}
			if got.HTTPStatus != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, got.HTTPStatus)
			}
			if !errors.Is(got.Err, originalErr) {
				t.Error("expected original error to be wrapped")
			}
			if got.Internal != internalInfo {
				t.Errorf("expected internal info %s, got %s", internalInfo, got.Internal)
			}
		})
	}
}

func TestNew_InternalArgs(t *testing.T) {
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
			name:         "Internal info provided",
			code:         ErrorCodeInternal,
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []string{internalInfo},
			wantInternal: internalInfo,
		},
		{
			name:         "No internal info",
			code:         ErrorCodeProxyCodeIssue,
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: nil,
			wantInternal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.code, tt.status, tt.err, tt.internalArgs...)

			if got.Internal != tt.wantInternal {
				t.Errorf("expected internal %q, got %q", tt.wantInternal, got.Internal)
			}
		})
	}
}
