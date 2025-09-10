package config

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPort(t *testing.T) {
	t.Run("with env var", func(t *testing.T) {
		t.Setenv("PORT", "12345")

		require.Equal(t, "12345", GetPort())
	})

	t.Run("without env var", func(t *testing.T) {
		t.Setenv("PORT", "")

		require.Equal(t, "9000", GetPort())
	})
}

func TestLoadFirebaseConfig(t *testing.T) {
	t.Run("loads from base64", func(t *testing.T) {
		dummy := "test-credentials"
		b64 := base64.StdEncoding.EncodeToString([]byte(dummy))
		t.Setenv("GOOGLE_APPLICATION_CREDENTIALS_BASE64", b64)

		cfg, err := LoadFirebaseConfig()
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Equal(t, []byte(dummy), cfg.CredentialsJSON)
	})

	t.Run("returns error if env is not set", func(t *testing.T) {
		t.Setenv("GOOGLE_APPLICATION_CREDENTIALS_BASE64", "")

		cfg, err := LoadFirebaseConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "GOOGLE_APPLICATION_CREDENTIALS_BASE64 is not set")
	})
}

func TestLoadGitHubOAuthConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("GITHUB_CLIENT_ID", "test-github-client-id")
		t.Setenv("GITHUB_REDIRECT_URI", "https://idpproxy.com/github/callback")

		cfg, err := LoadGitHubOAuthConfig()
		require.NoError(t, err)
		require.Equal(t, "test-github-client-id", cfg.ClientID)
		require.Equal(t, "https://idpproxy.com/github/callback", cfg.RedirectURI)
		require.Equal(t, "read:user", cfg.Scope)
		require.Equal(t, "true", cfg.AllowSignup)
	})

	t.Run("missing both variables", func(t *testing.T) {
		t.Setenv("GITHUB_CLIENT_ID", "")
		t.Setenv("GITHUB_REDIRECT_URI", "")

		cfg, err := LoadGitHubOAuthConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "GITHUB_CLIENT_ID and GITHUB_REDIRECT_URI are not set")
	})

	t.Run("missing client ID", func(t *testing.T) {
		t.Setenv("GITHUB_CLIENT_ID", "")
		t.Setenv("GITHUB_REDIRECT_URI", "https://idpproxy.com/github/callback")

		cfg, err := LoadGitHubOAuthConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "GITHUB_CLIENT_ID is not set")
	})

	t.Run("missing redirect URI", func(t *testing.T) {
		t.Setenv("GITHUB_CLIENT_ID", "test-github-client-id")
		t.Setenv("GITHUB_REDIRECT_URI", "")

		cfg, err := LoadGitHubOAuthConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "GITHUB_REDIRECT_URI is not set")
	})
}

func TestLoadGitHubDevConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("GITHUB_DEV_CLIENT_ID", "test-github-dev-client-id")
		t.Setenv("GITHUB_DEV_REDIRECT_URI", "http://localhost:9000/github/callback")

		cfg, err := LoadGitHubDevOAuthConfig()
		require.NoError(t, err)
		require.Equal(t, "test-github-dev-client-id", cfg.ClientID)
		require.Equal(t, "http://localhost:9000/github/callback", cfg.RedirectURI)
		require.Equal(t, "read:user", cfg.Scope)
		require.Equal(t, "true", cfg.AllowSignup)
	})

	t.Run("missing both variables", func(t *testing.T) {
		t.Setenv("GITHUB_DEV_CLIENT_ID", "")
		t.Setenv("GITHUB_DEV_REDIRECT_URI", "")

		cfg, err := LoadGitHubDevOAuthConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "GITHUB_DEV_CLIENT_ID and GITHUB_DEV_REDIRECT_URI are not set")
	})

	t.Run("missing client ID", func(t *testing.T) {
		t.Setenv("GITHUB_DEV_CLIENT_ID", "")
		t.Setenv("GITHUB_DEV_REDIRECT_URI", "http://localhost:9000/github/callback")

		cfg, err := LoadGitHubDevOAuthConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "GITHUB_DEV_CLIENT_ID is not set")
	})

	t.Run("missing redirect URI", func(t *testing.T) {
		t.Setenv("GITHUB_DEV_CLIENT_ID", "test-github-dev-client-id")
		t.Setenv("GITHUB_DEV_REDIRECT_URI", "")

		cfg, err := LoadGitHubDevOAuthConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "GITHUB_DEV_REDIRECT_URI is not set")
	})
}

func TestLoadServiceAccountConfig(t *testing.T) {
	t.Run("when env var is not set", func(t *testing.T) {
		t.Setenv("IMPERSONATE_SERVICE_ACCOUNT", "")
		cfg := LoadServiceAccountConfig()
		require.Equal(t, "", cfg.ImpersonateSA)
	})

	t.Run("when env var is set", func(t *testing.T) {
		t.Setenv("IMPERSONATE_SERVICE_ACCOUNT", "sa@example.iam.gserviceaccount.com")
		cfg := LoadServiceAccountConfig()
		require.Equal(t, "sa@example.iam.gserviceaccount.com", cfg.ImpersonateSA)
	})

	t.Run("when env var has spaces", func(t *testing.T) {
		t.Setenv("IMPERSONATE_SERVICE_ACCOUNT", "  sa@example.iam.gserviceaccount.com  ")
		cfg := LoadServiceAccountConfig()
		require.Equal(t, "sa@example.iam.gserviceaccount.com", cfg.ImpersonateSA)
	})
}
