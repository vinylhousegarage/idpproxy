package info_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestInfoRoute_Returns200AndJSONHealthy(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	di := testhelpers.NewMockSystemDeps(logger)

	r := router.NewRouter(di, http.FS(public.PublicFS))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/info", nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{
		"message": "Welcome to IdP Proxy",
		"openapi": "http://localhost:9000/openapi.json"
	}`, w.Body.String())
}
