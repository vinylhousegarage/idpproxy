package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"idpproxy/internal/router"

	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestHealthRoute_Returns200AndJSONHealthy(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	r := router.NewRouter(logger)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"status":"healthy"}`, w.Body.String())
}
