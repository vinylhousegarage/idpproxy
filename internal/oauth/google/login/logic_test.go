package login

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	const state = "abc123"
	c := BuildStateCookie(state)

	assert.Equal(t, "oauth_state", c.Name)
	assert.Equal(t, state, c.Value)
	assert.Equal(t, "/", c.Path)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
	assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
}

var mockCfg = &config.GoogleConfig{
	ClientID:     "client-id",
	ClientSecret: "client-secret",
	RedirectURI:  "https://localhost/callback",
	ResponseType: "code",
	Scope:        "openid email profile",
	AccessType:   "offline",
	Prompt:       "consent",
}

func TestBuildGoogleLoginURL_Success(t *testing.T) {
	t.Parallel()

	endpoint := "https://auth.example.com/oauth2/authorize"
	state := "sample-state-value"

	result, err := BuildGoogleLoginURL(mockCfg, endpoint, state)
	assert.NoError(t, err)

	parsed, err := url.Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, "auth.example.com", parsed.Host)
	assert.Equal(t, "/oauth2/authorize", parsed.Path)

	queries := parsed.Query()
	assert.Equal(t, "code", queries.Get("response_type"))
	assert.Equal(t, mockCfg.ClientID, queries.Get("client_id"))
	assert.Equal(t, mockCfg.RedirectURI, queries.Get("redirect_uri"))
	assert.Equal(t, mockCfg.Scope, queries.Get("scope"))
	assert.Equal(t, state, queries.Get("state"))
}

func TestBuildGoogleLoginURL_InvalidEndpoint(t *testing.T) {
	t.Parallel()

	_, err := BuildGoogleLoginURL(mockCfg, "://invalid-url", "state")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToParseLoginURL)
}
