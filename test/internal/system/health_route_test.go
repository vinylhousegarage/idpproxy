package system_test

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

func TestHealthRoute_Returns200AndJSONHealthy(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	githubAPIDeps := testhelpers.NewMockGitHubAPIDeps(logger)
	githubOAuthDeps := testhelpers.NewMockGitHubOAuthDeps(logger)
	googleDeps := testhelpers.NewMockGoogleDeps(logger)
	systemDeps := testhelpers.NewMockSystemDeps(logger)

	d := router.NewRouterDeps(public.PublicFS, githubAPIDeps, githubOAuthDeps, googleDeps, systemDeps)
	r := router.NewRouter(d)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"status":"healthy"}`, w.Body.String())
}
