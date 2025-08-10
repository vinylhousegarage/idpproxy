package user

import (
	"net/http"
	"testing"
)

func TestNewGitHubUserRequest_SetsMethodURLAndHeaders(t *testing.T) {
	t.Parallel()

	token := "testtoken123"

	req, err := NewGitHubUserRequest(token)
	if err != nil {
		t.Fatalf("NewGitHubUserRequest returned error: %v", err)
	}

	if req.Method != http.MethodGet {
		t.Errorf("method = %q, want %q", req.Method, http.MethodGet)
	}

	if got := req.URL.String(); got != GitHubUserURL {
		t.Errorf("url = %q, want %q", got, GitHubUserURL)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer "+token {
		t.Errorf("Authorization = %q, want %q", got, "Bearer "+token)
	}
	if got := req.Header.Get("Accept"); got != "application/vnd.github+json" {
		t.Errorf("Accept = %q, want %q", got, "application/vnd.github+json")
	}

	if req.Body != nil {
		t.Errorf("Body = non-nil, want nil")
	}
}

func TestNewGitHubUserRequest_EmptyTokenStillSetsBearerPrefix(t *testing.T) {
	t.Parallel()

	req, err := NewGitHubUserRequest("")
	if err != nil {
		t.Fatalf("NewGitHubUserRequest returned error: %v", err)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer " {
		t.Errorf("Authorization = %q, want %q", got, "Bearer ")
	}
}
