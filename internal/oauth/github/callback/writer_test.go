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
	rec := httptest.NewRecorder()

	logger := zap.NewNop()

	writeJSON(rec, http.StatusBadRequest, ErrorResponse{
		Error: apierror.ErrorMissingGitHubCode,
	}, logger)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("wrong content-type")
	}

	if rec.Body.Len() == 0 {
		t.Fatalf("body should not be empty")
	}

	var res ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != apierror.ErrorMissingGitHubCode {
		t.Fatalf("error = %s", res.Error)
	}
}
