package token

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func Exchange(
	ctx context.Context,
	httpc HTTPDoer,
	cfg *config.GitHubOAuthConfig,
	code string,
	state string,
) (string, error) {

	req, err := newTokenRequest(ctx, cfg, code, state)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	resp, err := httpc.Do(req)
	if err != nil {
		return "", fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%w: %d", ErrNon2xxStatus, resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	if len(raw) == 0 {
		return "", ErrEmptyBody
	}

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.Contains(ct, "application/json"):
		return parseJSON(raw)

	case strings.Contains(ct, "application/x-www-form-urlencoded"):
		return parseForm(raw)

	default:
		return "", fmt.Errorf("unsupported content-type: %s", ct)
	}
}

func newTokenRequest(
	ctx context.Context,
	cfg *config.GitHubOAuthConfig,
	code string,
	state string,
) (*http.Request, error) {

	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", cfg.RedirectURI)
	if state != "" {
		form.Set("state", state)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

type jsonResp struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

func parseJSON(b []byte) (string, error) {
	var r jsonResp
	if err := json.Unmarshal(b, &r); err != nil {
		return "", err
	}

	if r.Error != "" {
		return "", fmt.Errorf("%w: %s", ErrGitHubOAuthError, r.Error)
	}

	if r.AccessToken == "" {
		return "", ErrMissingAccessToken
	}

	return r.AccessToken, nil
}

func parseForm(b []byte) (string, error) {
	v, err := url.ParseQuery(string(b))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParseFormBody, err)
	}

	token := v.Get("access_token")
	if token == "" {
		return "", ErrMissingAccessToken
	}

	return token, nil
}
