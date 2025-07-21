package public_test

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

func TestLoginHTMLServed(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	di := testhelpers.NewMockDeps(logger)

	r := router.NewRouter(di, http.FS(public.PublicFS))

	req := httptest.NewRequest(http.MethodGet, "/static/login.html", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
	require.Contains(t, resp.Body.String(), "ログイン")
}
