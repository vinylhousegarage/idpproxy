package callback

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func BuildAccessTokenRequest(ctx context.Context, cfg *config.GitHubOAuthConfig, code, state string) (*http.Request, error) {
	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", cfg.RedirectURI)
	form.Set("state", state)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.GitHubTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	return req, nil
}
