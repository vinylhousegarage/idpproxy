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

func TestPrivacyPolicyPage(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	githubDeps := testhelpers.NewMockGitHubDeps(logger)
	githubAPIDeps := testhelpers.NewMockGitHubAPIDeps(logger)
	googleDeps := testhelpers.NewMockGoogleDeps(logger)
	systemDeps := testhelpers.NewMockSystemDeps(logger)

	r := router.NewRouter(githubDeps, githubAPIDeps, googleDeps, systemDeps, http.FS(public.PublicFS))

	req := httptest.NewRequest(http.MethodGet, "/privacy", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "プライバシーポリシー")
	require.Contains(t, resp.Body.String(), "<!DOCTYPE html>")
}
