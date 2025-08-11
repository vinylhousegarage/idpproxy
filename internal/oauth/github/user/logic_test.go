package user

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
)

func TestNewGitHubUserRequest_SetsMethodURLAndHeaders(t *testing.T) {
	t.Parallel()

	token := "testtoken123"

	req, err := NewGitHubUserRequest(context.Background(), token)
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, githubUserURL, req.URL.String())
	require.Equal(t, "Bearer "+token, req.Header.Get("Authorization"))
	require.Equal(t, "application/vnd.github+json", req.Header.Get("Accept"))
	require.Equal(t, config.GitHubAPIVersion, req.Header.Get("X-GitHub-Api-Version"))
	require.Equal(t, config.UserAgent(), req.Header.Get("User-Agent"))
	require.Nil(t, req.Body)
}

func TestNewGitHubUserRequest_EmptyTokenReturnsError(t *testing.T) {
	t.Parallel()

	_, err := NewGitHubUserRequest(context.Background(), "")
	require.ErrorIs(t, err, ErrEmptyBearerToken)
}

//nolint:staticcheck // SA1012: intentionally passing nil to verify ErrNilContext
func TestNewGitHubUserRequest_NilContextReturnsError(t *testing.T) {
	t.Parallel()

	_, err := NewGitHubUserRequest(nil, "token")
	require.ErrorIs(t, err, ErrNilContext)
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
		resp := newHTTPResponse(t, http.StatusOK, body, &closed)

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
		resp := newHTTPResponse(t, http.StatusUnauthorized, `{"message":"bad creds"}`, &closed)

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
		resp := newHTTPResponse(t, http.StatusOK, `{"id": 1, "login":`, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, closed)
		require.Contains(t, err.Error(), "failed to decode GitHub user response")
	})

	t.Run("OmittedFields", func(t *testing.T) {
		t.Parallel()

		closed := false
		resp := newHTTPResponse(t, http.StatusOK, `{"id":123,"login":"octocat"}`, &closed)

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
		resp := newHTTPResponse(t, http.StatusOK, ``, &closed)

		got, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, closed)
		require.Contains(t, err.Error(), "failed to decode GitHub user response")
	})

	t.Run("Non2xx_SnippetIsTruncated", func(t *testing.T) {
		t.Parallel()

		head := bytes.Repeat([]byte{'x'}, 256)
		tail := bytes.Repeat([]byte{'Y'}, 44)
		payload := append(head, tail...)

		closed := false
		resp := newHTTPResponse(t, http.StatusForbidden, string(payload), &closed)

		_, err := DecodeGitHubUserResponse(resp)
		require.Error(t, err)
		require.True(t, closed)

		msg := err.Error()
		require.Contains(t, msg, "non-2xx status")
		require.Contains(t, msg, "403")

		require.Contains(t, msg, string(head))
		require.NotContains(t, msg, string(tail))
	})
}
