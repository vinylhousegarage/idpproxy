package token

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

func TestBuildAccessTokenRequest(t *testing.T) {
	t.Parallel()

	cfg := &config.GitHubOAuthConfig{
		ClientID:     "client-id-123",
		ClientSecret: "secret-xyz",
		RedirectURI:  "https://idpproxy.com/github/callback",
	}

	readBody := func(t *testing.T, r io.ReadCloser) []byte {
		t.Helper()
		b, err := io.ReadAll(r)
		require.NoError(t, err)
		require.NoError(t, r.Close())
		return b
	}

	t.Run("builds POST request with correct URL and headers", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		req, err := BuildAccessTokenRequest(ctx, cfg, "abc-code", "state-1")
		require.NoError(t, err)

		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, config.GitHubTokenURL, req.URL.String())

		require.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))

		body := readBody(t, req.Body)
		require.Greater(t, len(body), 0)
	})

	t.Run("encodes form body correctly", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		code := "gho_CODE123"
		state := "randomSTATE"

		req, err := BuildAccessTokenRequest(ctx, cfg, code, state)
		require.NoError(t, err)

		body := readBody(t, req.Body)

		values, err := url.ParseQuery(string(body))
		require.NoError(t, err)

		require.Equal(t, cfg.ClientID, values.Get("client_id"))
		require.Equal(t, cfg.ClientSecret, values.Get("client_secret"))
		require.Equal(t, code, values.Get("code"))
		require.Equal(t, cfg.RedirectURI, values.Get("redirect_uri"))
		require.Equal(t, state, values.Get("state"))

		require.Len(t, values, 5)
	})

	t.Run("redirect_uri is safely url-encoded in body", func(t *testing.T) {
		t.Parallel()

		cfg2 := *cfg
		cfg2.RedirectURI = "https://idpproxy.com/github/callback?next=/welcome page&ref=gh"

		req, err := BuildAccessTokenRequest(context.Background(), &cfg2, "code-1", "state-2")
		require.NoError(t, err)

		body := readBody(t, req.Body)

		values, err := url.ParseQuery(string(body))
		require.NoError(t, err)
		require.Equal(t, cfg2.RedirectURI, values.Get("redirect_uri"))

		raw := string(body)
		require.Contains(t, raw, "redirect_uri=")
		require.NotContains(t, raw, "welcome page")
	})

	t.Run("request is bound to provided context", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		req, err := BuildAccessTokenRequest(ctx, cfg, "code-ctx", "state-ctx")
		require.NoError(t, err)

		cancel()

		select {
		case <-req.Context().Done():
			require.ErrorIs(t, req.Context().Err(), context.Canceled)
		case <-time.After(200 * time.Millisecond):
			t.Fatal("request context was not canceled")
		}
	})
}
