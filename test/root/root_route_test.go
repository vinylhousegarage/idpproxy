package root_test

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

func TestRootServesLoginHTML(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	githubDeps := testhelpers.NewMockGitHubDeps(logger)
	googleDeps := testhelpers.NewMockGoogleDeps(logger)
	systemDeps := testhelpers.NewMockSystemDeps(logger)

	r := router.NewRouter(githubDeps, googleDeps, systemDeps, http.FS(public.PublicFS))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
	require.Contains(t, resp.Body.String(), "ログイン")
}
