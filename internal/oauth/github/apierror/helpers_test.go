package apierror

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestAPIErrors(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")

	tests := []struct {
		name           string
		fn             func(error, ...string) *APIError
		expectedCode   ErrorCode
		expectedStatus int
		internalInfo   []string
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
		{
			name:           "GitHubAccessTokenRequestError",
			fn:             GitHubAccessTokenRequestError,
			expectedCode:   ErrorCodeGitHubAccessTokenRequest,
			expectedStatus: http.StatusBadGateway,
		},
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

			got := tt.fn(originalErr, tt.internalInfo...)

			if got.Code != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, got.Code)
			}
			if got.HTTPStatus != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, got.HTTPStatus)
			}
			if !errors.Is(got.Err, originalErr) {
				t.Error("expected original error to be wrapped")
			}
			if !reflect.DeepEqual(got.Internal, tt.internalInfo) {
				t.Errorf("expected internal %v, got %v", tt.internalInfo, got.Internal)
			}
		})
	}
}
