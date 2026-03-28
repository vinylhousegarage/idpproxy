package callback

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"
	"go.uber.org/zap"
)

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	logger := zap.NewNop()

	expectedStatus := http.StatusBadRequest
	expectedError := string(apierror.ErrorMissingGitHubCode)
	expectedContentType := "application/json; charset=utf-8"

	writeJSON(rec, expectedStatus, apierror.ErrorResponse{
		Error: expectedError,
	}, logger)

	if rec.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, rec.Code)
	}

	if rec.Header().Get("Content-Type") != expectedContentType {
		t.Fatalf("expected content-type %s, got %s",
			expectedContentType,
			rec.Header().Get("Content-Type"),
		)
	}

	if rec.Body.Len() == 0 {
		t.Fatalf("response body should not be empty")
	}

	var res apierror.ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != expectedError {
		t.Fatalf("expected error %s, got %s", expectedError, res.Error)
	}
}
