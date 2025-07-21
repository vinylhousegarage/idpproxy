package public_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/router"
)

func TestLoginHTMLServed(t *testing.T) {
	t.Parallel()

	di := deps.NewTestDependencies(t)
	r := router.NewRouter(di)

	req := httptest.NewRequest(http.MethodGet, "/public/login.html", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
	require.Contains(t, resp.Body.String(), "ログイン")
}
