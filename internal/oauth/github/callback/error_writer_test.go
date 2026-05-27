package callback

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

func TestWriteError_WithAPIError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	err := apierror.New(apierror.ErrorCodeMissingState, http.StatusBadRequest, errors.New("missing state"))

	WriteError(c, err)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var res apierror.ErrorResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != apierror.ErrorCodeMissingState {
		t.Fatalf("expected %s, got %s", apierror.ErrorCodeMissingState, res.Error)
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

	if res.Error != apierror.ErrorCodeInternal {
		t.Fatalf("expected %s, got %s", apierror.ErrorCodeInternal, res.Error)
	}
}
