package root_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/router"

	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestRootRoute_Returns200AndJSONHealthy(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}
	r := router.NewRouter(logger)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{
		"message":"Welcome to IdP Proxy",
		"openapi": "http://localhost:9000/openapi.json",
	}`, w.Body.String())
}
