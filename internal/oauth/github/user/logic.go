package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
)

var (
	ErrEmptyBearerToken                 = errors.New("empty bearer token")
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid Authorization header format")
	ErrMissingAuthorizationHeader       = errors.New("missing Authorization header")
	ErrNilContext                       = errors.New("nil context")
)

var githubUserURL = config.GitHubUserURL

func ExtractAuthHeaderToken(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", ErrMissingAuthorizationHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidAuthorizationHeaderFormat
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrEmptyBearerToken
	}
	return token, nil
}

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
	req.Header.Set("X-GitHub-Api-Version", config.GitHubAPIVersion)
	req.Header.Set("User-Agent", config.UserAgent())

	return req, nil
}

func DecodeGitHubUserResponse(resp *http.Response) (*response.GitHubUserAPIResponse, error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return nil, fmt.Errorf("non-2xx status: %d body=%q", resp.StatusCode, snippet)
	}

	var githubUserAPIResponse response.GitHubUserAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&githubUserAPIResponse); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub user response: %w", err)
	}

	return &githubUserAPIResponse, nil
}
