package user

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

var (
	ErrNilContext       = errors.New("nil context")
	ErrEmptyBearerToken = errors.New("empty bearer token")
)

var githubUserURL = config.GitHubUserURL

func NewGitHubUserRequest(ctx context.Context, accessToken string) (*http.Request, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	token := strings.TrimSpace(accessToken)
	if token == "" {
		return nil, ErrEmptyBearerToken
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	return req, nil
}
