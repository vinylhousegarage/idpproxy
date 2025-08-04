package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestGitHubLoginHandler_Serve(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	dependencies := testhelpers.NewMockGitHubDeps(logger)
	handler := NewGitHubLoginHandler(dependencies)

	router := http.NewServeMux()
	router.HandleFunc("/github/login", func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		handler.Serve(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/github/login")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	require.NotEmpty(t, location)
	require.True(t, strings.HasPrefix(location, "https://github.com/login/oauth/authorize?"))
	require.Contains(t, location, "client_id=test-client-id")
	require.Contains(t, location, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
	require.Contains(t, location, "scope=read%3Auser")
	require.Contains(t, location, "allow_signup=false")
	require.Contains(t, location, "state=")

	cookies := resp.Cookies()
	foundStateCookie := false
	for _, c := range cookies {
		if c.Name == "oauth_state" && c.HttpOnly && c.Secure {
			foundStateCookie = true
		}
	}
	require.True(t, foundStateCookie, "oauth_state cookie should be set and secure")
}
