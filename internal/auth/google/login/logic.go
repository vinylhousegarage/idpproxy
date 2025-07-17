package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func GenerateState() string {
	const stateLength = 16
	b := make([]byte, stateLength)
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

func BuildGoogleLoginURL(cfg config.GoogleConfig, endpoint, state string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", ErrFailedToParseLoginURL
	}
	q := u.Query()
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURI)
	q.Set("response_type", cfg.ResponseType)
	q.Set("scope", cfg.Scope)
	q.Set("state", state)
	q.Set("access_type", cfg.AccessType)
	q.Set("prompt", cfg.Prompt)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
