package public_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

//go:embed public/*
var publicFS embed.FS

func TestLoginHTMLServed(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	di := testhelpers.NewMockDeps(logger)

	r := router.NewRouter(di, http.FS(publicFS))

	req := httptest.NewRequest(http.MethodGet, "/public/login.html", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
	require.Contains(t, resp.Body.String(), "ログイン")
}
