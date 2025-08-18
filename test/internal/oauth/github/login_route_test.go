package github_test

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

func TestGitHubLoginRoute_Returns302Redirect(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	githubDeps := testhelpers.NewMockGitHubDeps(logger)
	googleDeps := testhelpers.NewMockGoogleDeps(logger)
	systemDeps := testhelpers.NewMockSystemDeps(logger)

	r := router.NewRouter(githubDeps, githubAPIDeps, googleDeps, systemDeps, http.FS(public.PublicFS))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/github/login", nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusFound, w.Code)

	location := w.Header().Get("Location")
	require.NotEmpty(t, location)
	require.Contains(t, location, "https://github.com/login/oauth/authorize?")
	require.Contains(t, location, "client_id=test-client-id")
	require.Contains(t, location, "state=")

	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "oauth_state" {
			found = true
			break
		}
	}
	require.True(t, found, "oauth_state cookie should be set")
}
