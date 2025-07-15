package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"yourmodule/router"
)

func TestHealthRouteReturns_JSONHealthy(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	r := router.NewRouter(logger)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	expected := `{"status":"healthy"}`
	if w.Body.String() != expected {
		t.Fatalf("unexpected body: got %s, want %s", w.Body.String(), expected)
	}
}
