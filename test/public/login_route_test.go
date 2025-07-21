package public_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat("/app"); err == nil {
		_ = os.Chdir("/app")
	}

	os.Exit(m.Run())
}

func TestLoginHTMLServed(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	di := testhelpers.NewMockDeps(logger)
	r := router.NewRouter(di)

	req := httptest.NewRequest(http.MethodGet, "/public/login.html", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
	require.Contains(t, resp.Body.String(), "ログイン")
}
