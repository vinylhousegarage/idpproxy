package user

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/response"
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

	if _, err := NewGitHubUserRequest(context.Background(), ""); !errors.Is(err, ErrEmptyBearerToken) {
		t.Fatalf("want ErrEmptyBearerToken, got %v", err)
	}
}

//nolint:staticcheck // SA1012: intentionally passing nil to verify ErrNilContext
func TestNewGitHubUserRequest_NilContextReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := NewGitHubUserRequest(nil, "token"); !errors.Is(err, ErrNilContext) {
		t.Fatalf("want ErrNilContext, got %v", err)
	}
}

type testReadCloser struct {
	r      io.Reader
	closed *bool
}

func (t *testReadCloser) Read(p []byte) (int, error) { return t.r.Read(p) }
func (t *testReadCloser) Close() error {
	if t.closed != nil {
		*t.closed = true
	}
	return nil
}

func newHTTPResponse(t *testing.T, status int, body string, closed *bool) *http.Response {
	t.Helper()

	return &http.Response{
		StatusCode: status,
		Body:       &testReadCloser{r: bytes.NewBufferString(body), closed: closed},
		Header:     make(http.Header),
	}
}

func TestDecodeGitHubUserResponse(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		body := `{"id":123,"login":"octocat","email":"octo@example.com","name":"The Octocat"}`
		closed := false
		resp := newHTTPResponse(http.StatusOK, body, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.NoError(t, err)
		require.True(t, closed)

		want := &response.GitHubUserAPIResponse{
			ID:    123,
			Login: "octocat",
			Email: "octo@example.com",
			Name:  "The Octocat",
		}
		require.Equal(t, want, got)
	})

	t.Run("Non2xx", func(t *testing.T) {
		t.Parallel()

		closed := false
		resp := newHTTPResponse(http.StatusUnauthorized, `{"message":"bad creds"}`, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, closed)
		require.Contains(t, err.Error(), "non-2xx status")
		require.Contains(t, err.Error(), "401")
		require.Contains(t, err.Error(), "bad creds")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		t.Parallel()

		closed := false
		resp := newHTTPResponse(http.StatusOK, `{"id": 1, "login":`, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, closed)
		require.Contains(t, err.Error(), "failed to decode GitHub user response")
	})

	t.Run("OmittedFields", func(t *testing.T) {
		t.Parallel()

		closed := false
		resp := newHTTPResponse(http.StatusOK, `{"id":123,"login":"octocat"}`, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.NoError(t, err)
		require.True(t, closed)
		require.Equal(t, int64(123), got.ID)
		require.Equal(t, "octocat", got.Login)
		require.Empty(t, got.Email)
		require.Empty(t, got.Name)
	})

	t.Run("EmptyBody", func(t *testing.T) {
		t.Parallel()

		closed := false
		resp := newHTTPResponse(http.StatusOK, ``, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, closed)
		require.Contains(t, err.Error(), "failed to decode GitHub user response")
	})

	t.Run("Non2xx_SnippetIsTruncated", func(t *testing.T) {
		t.Parallel()

		long := make([]byte, 300)
		for i := range long {
			long[i] = 'x'
		}
		closed := false
		resp := newHTTPResponse(http.StatusForbidden, string(long), &closed)

		_, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.True(t, closed)

		msg := err.Error()
		require.Contains(t, msg, "non-2xx status")
		require.Contains(t, msg, "403")

		require.Contains(t, msg, string(long[:256]))
		require.NotContains(t, msg, string(long[256:]))
	})
}
