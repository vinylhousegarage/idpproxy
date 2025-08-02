package login

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildGitHubLoginURL(t *testing.T) {
	clientID := "dummy-client-id"
	redirectURI := "https://idpproxy.com/github/callback"

	expected := "https://github.com/login/oauth/authorize" +
		"?client_id=" + url.QueryEscape(clientID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&scope=read:user+user:email"

	actual := BuildGitHubLoginURL(clientID, redirectURI)

	require.Equal(t, expected, actual)
}
