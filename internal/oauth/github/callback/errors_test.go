package callback

import (
	"errors"
	"net/http"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

func TestCallbackErrors(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")
	internalInfo := "debug info"

	tests := []struct {
		name           string
		fn             func(error, ...string) *apierror.APIError
		expectedCode   apierror.ErrorCode
		expectedStatus int
	}{
		{
			name:           "MissingGitHubCode",
			fn:             MissingGitHubCode,
			expectedCode:   ErrorCodeMissingGitHubCode,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "MissingState",
			fn:             MissingState,
			expectedCode:   ErrorCodeMissingState,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidState",
			fn:             InvalidState,
			expectedCode:   ErrorCodeInvalidState,
			expectedStatus: http.StatusBadRequest,
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
