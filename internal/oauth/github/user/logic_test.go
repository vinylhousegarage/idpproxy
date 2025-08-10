package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func TestNewGitHubUserRequest_SetsMethodURLAndHeaders(t *testing.T) {
	t.Parallel()

	token := "testtoken123"

	req, err := NewGitHubUserRequest(context.Background(), token)
	if err != nil {
		t.Fatalf("NewGitHubUserRequest returned error: %v", err)
	}

	if req.Method != http.MethodGet {
		t.Errorf("method = %q, want %q", req.Method, http.MethodGet)
	}

	if got := req.URL.String(); got != githubUserURL {
		t.Errorf("url = %q, want %q", got, githubUserURL)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer "+token {
		t.Errorf("Authorization = %q, want %q", got, "Bearer "+token)
	}
	if got := req.Header.Get("Accept"); got != "application/vnd.github+json" {
		t.Errorf("Accept = %q, want %q", got, "application/vnd.github+json")
	}
	if got := req.Header.Get("X-GitHub-Api-Version"); got != config.GitHubAPIVersion {
		t.Errorf("X-GitHub-Api-Version = %q, want %q", got, config.GitHubAPIVersion)
	}
	if got := req.Header.Get("User-Agent"); got != config.UserAgent() {
		t.Errorf("User-Agent = %q, want %q", got, config.UserAgent())
	}

	if req.Body != nil {
		t.Errorf("Body = non-nil, want nil")
	}
}

func TestNewGitHubUserRequest_EmptyTokenReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := NewGitHubUserRequest(context.Background(), ""); err != ErrEmptyBearerToken {
		t.Fatalf("want ErrEmptyBearerToken, got %v", err)
	}
}

func TestNewGitHubUserRequest_NilContextReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := NewGitHubUserRequest(nil, "token"); err != ErrNilContext {
		t.Fatalf("want ErrNilContext, got %v", err)
	}
}
