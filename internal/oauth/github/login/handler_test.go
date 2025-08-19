package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestGitHubLoginHandler_Serve(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	dependencies := testhelpers.NewMockGitHubOAuthDeps(logger)
	handler := NewGitHubLoginHandler(dependencies)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/github/login", nil)

	handler.Serve(c)
	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusFound, res.StatusCode)

	location := res.Header.Get("Location")
	require.NotEmpty(t, location)
	require.True(t, strings.HasPrefix(location, "https://github.com/login/oauth/authorize?"))
	require.Contains(t, location, "client_id=test-client-id")
	require.Contains(t, location, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
	require.Contains(t, location, "scope=read%3Auser")
	require.Contains(t, location, "allow_signup=false")
	require.Contains(t, location, "state=")

	cookies := res.Cookies()
	foundStateCookie := false
	for _, c := range cookies {
		if c.Name == "oauth_state" && c.HttpOnly && c.Secure {
			foundStateCookie = true
			break
		}
	}
	require.True(t, foundStateCookie, "oauth_state cookie should be set and secure")
}
