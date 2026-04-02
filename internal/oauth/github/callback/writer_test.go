package callback

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"
)

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	writeJSON(rec, http.StatusBadRequest, apierror.ErrorResponse{
		Error: string(apierror.ErrorMissingGitHubCode),
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("wrong content-type")
	}

	if rec.Body.Len() == 0 {
		t.Fatalf("body should not be empty")
	}

	var res apierror.ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != string(apierror.ErrorMissingGitHubCode) {
		t.Fatalf("error = %s", res.Error)
	}
}

func TestWriteError_WithAPIError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	err := apierror.New(apierror.ErrorMissingState, http.StatusBadRequest, errors.New("missing state"))

	WriteError(c, err)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var res apierror.ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != string(apierror.ErrorMissingState) {
		t.Fatalf("expected %s, got %s", apierror.ErrorMissingState, res.Error)
	}
}

func TestWriteError_WithGenericError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	err := errors.New("unexpected error")

	WriteError(c, err)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var res apierror.ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != string(apierror.ErrorInternal) {
		t.Fatalf("expected %s, got %s", apierror.ErrorInternal, res.Error)
	}
}
