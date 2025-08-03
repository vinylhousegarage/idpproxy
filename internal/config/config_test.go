package config

import (
	"encoding/base64"
	"os"
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

func TestLoadGitHubConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("GITHUB_CLIENT_ID", "test-github-client-id")
		t.Setenv("GITHUB_CLIENT_SECRET", "test-github-client-secret")
		t.Setenv("GITHUB_REDIRECT_URI", "http://localhost:9000/github/callback")

		cfg, err := LoadGitHubConfig()
		require.NoError(t, err)
		require.Equal(t, "test-github-client-id", cfg.ClientID)
		require.Equal(t, "test-github-client-secret", cfg.ClientSecret)
		require.Equal(t, "http://localhost:9000/github/callback", cfg.RedirectURI)
		require.Equal(t, "read:user user:email", cfg.Scope)
		require.Equal(t, "true", cfg.AllowSignup)
	})

	t.Run("missing required variables", func(t *testing.T) {
		t.Setenv("GITHUB_CLIENT_ID", "")
		t.Setenv("GITHUB_CLIENT_SECRET", "")
		t.Setenv("GITHUB_REDIRECT_URI", "")

		cfg, err := LoadGitHubConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing required environment variables")
	})
}
