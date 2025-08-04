package login

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func TestGenerateState(t *testing.T) {
	t.Parallel()

	state1 := GenerateState()
	state2 := GenerateState()

	require.NotEmpty(t, state1)
	require.NotEmpty(t, state2)
	require.NotEqual(t, state1, state2, "state should be random")
}

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	state := "teststate"
	cookie := BuildStateCookie(state)

	require.Equal(t, "oauth_state", cookie.Name)
	require.Equal(t, state, cookie.Value)
	require.True(t, cookie.HttpOnly)
	require.True(t, cookie.Secure)
	require.Equal(t, "/", cookie.Path)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
}

func TestBuildGitHubLoginURL(t *testing.T) {
	t.Parallel()

	cfg := &config.GitHubConfig{
		ClientID:    "test-client-id",
		RedirectURI: "https://example.com/callback",
		Scope:       "read:user",
		AllowSignup: "false",
	}
	state := "teststate"

	urlStr := BuildGitHubLoginURL(cfg, state)
	require.True(t, strings.HasPrefix(urlStr, "https://github.com/login/oauth/authorize?"))

	u, err := url.Parse(urlStr)
	require.NoError(t, err)

	q := u.Query()
	require.Equal(t, cfg.ClientID, q.Get("client_id"))
	require.Equal(t, cfg.RedirectURI, q.Get("redirect_uri"))
	require.Equal(t, cfg.Scope, q.Get("scope"))
	require.Equal(t, state, q.Get("state"))
	require.Equal(t, cfg.AllowSignup, q.Get("allow_signup"))
}
