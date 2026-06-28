package apierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/gin-gonic/gin"
)

var logger = zap.NewNop()

func TestErrorLogger_WithAPIError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ErrorLogger(logger))

	r.GET("/test", func(c *gin.Context) {
		err := New(ErrorCodeMissingState, http.StatusBadRequest, errors.New("missing state"))
		_ = c.Error(err)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var res ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != ErrorCodeMissingState {
		t.Fatalf("expected %s, got %s", ErrorCodeMissingState, res.Error)
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

	var res ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != ErrorCodeInternal {
		t.Fatalf("expected %s, got %s", ErrorCodeInternal, res.Error)
	}
}

func TestErrorLogger_WithWrappedAPIError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ErrorLogger(logger))

	r.GET("/test", func(c *gin.Context) {
		apiErr := New(ErrorCodeMissingState, http.StatusBadRequest, errors.New("missing state"))
		wrapped := fmt.Errorf("wrapped: %w", apiErr)
		_ = c.Error(wrapped)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var res ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if res.Error != ErrorCodeMissingState {
		t.Fatalf("expected %s, got %s", ErrorCodeMissingState, res.Error)
	}
}

func TestErrorLogger_WithStatus400Error_LogsCorrectFields(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.WarnLevel)
	observedLogger := zap.New(core)

	r := gin.New()
	r.Use(ErrorLogger(observedLogger))

	r.GET("/test", func(c *gin.Context) {
		apiErr := New(ErrorCodeMissingState, http.StatusBadRequest, errors.New("missing state"))
		apiErr.Internal = []string{"debug details here"}
		_ = c.Error(apiErr)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	logEntry := logs.All()[0]

	if logEntry.Message != "request failed" {
		t.Errorf("expected message 'request failed', got '%s'", logEntry.Message)
	}

	fields := logEntry.ContextMap()

	expectedFields := map[string]interface{}{
		"path":     "/test",
		"method":   "GET",
		"code":     string(ErrorCodeMissingState),
		"status":   int64(http.StatusBadRequest),
		"detail_1": "debug details here",
	}

	for k, expectedVal := range expectedFields {
		if gotVal, ok := fields[k]; !ok {
			t.Errorf("expected log field '%s' to be present", k)
		} else if gotVal != expectedVal {
			t.Errorf("expected log field '%s' to be %v, got %v", k, expectedVal, gotVal)
		}
	}

	if errField, ok := fields["error"]; !ok {
		t.Errorf("expected 'error' field in log")
	} else if gotErrStr, ok := errField.(string); !ok {
		t.Errorf("expected 'error' field to be a string")
	} else if gotErrStr != "missing state" {
		t.Errorf("expected log error message 'missing state', got '%s'", gotErrStr)
	}
}

func TestErrorLogger_WithStatus500Error_LogsCorrectFields(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.ErrorLevel)
	observedLogger := zap.New(core)

	r := gin.New()
	r.Use(ErrorLogger(observedLogger))

	r.GET("/test", func(c *gin.Context) {
		apiErr := New(ErrorCodeInternal, http.StatusInternalServerError, errors.New("internal error"))
		apiErr.Internal = []string{"debug details here"}
		_ = c.Error(apiErr)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	logEntry := logs.All()[0]

	if logEntry.Message != "request failed" {
		t.Errorf("expected message 'request failed', got '%s'", logEntry.Message)
	}

	fields := logEntry.ContextMap()

	expectedFields := map[string]interface{}{
		"path":     "/test",
		"method":   "GET",
		"code":     string(ErrorCodeInternal),
		"status":   int64(http.StatusInternalServerError),
		"detail_1": "debug details here",
	}

	for k, expectedVal := range expectedFields {
		gotVal, ok := fields[k]
		if !ok {
			t.Errorf("expected log field '%s' to be present", k)
			continue
		}

		if reflect.TypeOf(expectedVal).Kind() == reflect.Slice {
			if !reflect.DeepEqual(gotVal, expectedVal) {
				t.Errorf("expected log field '%s' to be %v, got %v", k, expectedVal, gotVal)
			}
		} else {
			if gotVal != expectedVal {
				t.Errorf("expected log field '%s' to be %v, got %v", k, expectedVal, gotVal)
			}
		}
	}

	if errField, ok := fields["error"]; !ok {
		t.Errorf("expected 'error' field in log")
	} else if gotErrStr, ok := errField.(string); !ok {
		t.Errorf("expected 'error' field to be a string")
	} else if gotErrStr != "internal error" {
		t.Errorf("expected log error message 'internal error', got '%s'", gotErrStr)
	}
}

func TestErrorLogger_WithMultipleInternalInfo_LogsCorrectFields(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	core, logs := observer.New(zap.ErrorLevel)
	observedLogger := zap.New(core)

	r := gin.New()
	r.Use(ErrorLogger(observedLogger))

	r.GET("/test", func(c *gin.Context) {
		apiErr := New(ErrorCodeInternal, http.StatusInternalServerError, errors.New("multiple info error"))

		apiErr.Internal = []string{"first debug info", "second debug info"}
		_ = c.Error(apiErr)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	fields := logs.All()[0].ContextMap()

	expected := map[string]interface{}{
		"detail_1": "first debug info",
		"detail_2": "second debug info",
	}

	for k, v := range expected {
		if got, ok := fields[k]; !ok {
			t.Errorf("expected log field '%s' to be present", k)
		} else if got != v {
			t.Errorf("expected log field '%s' to be %v, got %v", k, v, got)
		}
	}
}
