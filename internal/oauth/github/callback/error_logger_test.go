package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"
)

var logger = zap.NewNop()

func TestErrorLogger_WithAPIError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ErrorLogger(logger))

	r.GET("/test", func(c *gin.Context) {
		err := apierror.New(apierror.ErrorMissingState, http.StatusBadRequest, errors.New("missing state"))
		_ = c.Error(err)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var res apierror.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != apierror.ErrorMissingState {
		t.Fatalf("expected %s, got %s", apierror.ErrorMissingState, res.Error)
	}
}

func TestErrorLogger_WithGenericError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ErrorLogger(logger))

	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("unexpected error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var res apierror.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != apierror.ErrorInternal {
		t.Fatalf("expected %s, got %s", apierror.ErrorInternal, res.Error)
	}
}

func TestErrorLogger_WithWrappedAPIError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ErrorLogger(logger))

	r.GET("/test", func(c *gin.Context) {
		apiErr := apierror.New(apierror.ErrorMissingState, http.StatusBadRequest, errors.New("missing state"))
		wrapped := fmt.Errorf("wrapped: %w", apiErr)
		_ = c.Error(wrapped)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var res apierror.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != apierror.ErrorMissingState {
		t.Fatalf("expected %s, got %s", apierror.ErrorMissingState, res.Error)
	}
}
