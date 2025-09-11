package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func GenerateState() string {
	b := make([]byte, config.OAuthStateBytes)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure random state: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}

func BuildStateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func BuildGitHubLoginURL(cfg *config.GitHubOAuthConfig, state string) string {
	v := url.Values{}
	v.Set("client_id", cfg.ClientID)
	v.Set("redirect_uri", cfg.RedirectURI)
	v.Set("scope", cfg.Scope)
	v.Set("state", state)
	v.Set("allow_signup", cfg.AllowSignup)

	return config.GitHubAuthorizeURL + "?" + v.Encode()
}
