package apierror

import (
	"errors"
	"testing"
)

func TestFromInternal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputErr       error
		expectedCode   ErrorCode
		expectedStatus int
	}{
		{
			name:           "invalid input",
			inputErr:       ErrInvalidInput,
			expectedCode:   ErrorInvalidState,
			expectedStatus: 400,
		},
		{
			name:           "persist error",
			inputErr:       ErrPersist,
			expectedCode:   ErrorProxyCodeIssue,
			expectedStatus: 500,
		},
		{
			name:           "unknown error",
			inputErr:       errors.New("something went wrong"),
			expectedCode:   ErrorInternal,
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FromInternal(tt.inputErr)

			if got == nil {
				t.Fatalf("expected APIError, got nil")
			}

			if got.Code != tt.expectedCode {
				t.Fatalf("code = %v, want %v", got.Code, tt.expectedCode)
			}

			if got.HTTPStatus != tt.expectedStatus {
				t.Fatalf("status = %v, want %v", got.HTTPStatus, tt.expectedStatus)
			}

			if !errors.Is(got.Err, tt.inputErr) {
				t.Fatalf("wrapped error mismatch")
			}
		})
	}
}
