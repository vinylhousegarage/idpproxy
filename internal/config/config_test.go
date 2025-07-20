package config

import (
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
		original := os.Getenv("PORT")
		os.Unsetenv("PORT")
		defer os.Setenv("PORT", original)

		require.Equal(t, "9000", GetPort())
	})
}

func TestGetOpenAPIURL(t *testing.T) {
	t.Setenv("OPENAPI_URL", "https://test.example.com/openapi.json")

	require.Equal(t, "https://test.example.com/openapi.json", GetOpenAPIURL())
}

func TestLoadGoogleConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
		t.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
		t.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:9000/callback")
		t.Setenv("GOOGLE_RESPONSE_TYPE", "code")
		t.Setenv("GOOGLE_SCOPE", "unused")
		t.Setenv("GOOGLE_ACCESS_TYPE", "unused")
		t.Setenv("GOOGLE_PROMPT", "unused")

		cfg, err := LoadGoogleConfig()
		require.NoError(t, err)
		require.Equal(t, "test-client-id", cfg.ClientID)
		require.Equal(t, "test-client-secret", cfg.ClientSecret)
		require.Equal(t, "http://localhost:9000/callback", cfg.RedirectURI)
		require.Equal(t, "code", cfg.ResponseType)
		require.Equal(t, "openid email profile", cfg.Scope)
		require.Equal(t, "offline", cfg.AccessType)
		require.Equal(t, "consent", cfg.Prompt)
	})

	t.Run("missing required variables", func(t *testing.T) {
		originalID := os.Getenv("GOOGLE_CLIENT_ID")
		originalSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
		originalRedirect := os.Getenv("GOOGLE_REDIRECT_URI")
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
		os.Unsetenv("GOOGLE_REDIRECT_URI")
		defer func() {
			os.Setenv("GOOGLE_CLIENT_ID", originalID)
			os.Setenv("GOOGLE_CLIENT_SECRET", originalSecret)
			os.Setenv("GOOGLE_REDIRECT_URI", originalRedirect)
		}()

		cfg, err := LoadGoogleConfig()
		require.Nil(t, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing required environment variables")
	})
}

func TestLoadFernetKeyString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("FERNET_KEY", "dummy_key")

		key, err := LoadFernetKeyString()
		require.NoError(t, err)
		require.Equal(t, "dummy_key", key)
	})

	t.Run("missing", func(t *testing.T) {
		original := os.Getenv("FERNET_KEY")
		os.Unsetenv("FERNET_KEY")
		defer os.Setenv("FERNET_KEY", original)

		_, err := LoadFernetKeyString()
		require.Error(t, err)
		require.Contains(t, err.Error(), "FERNET_KEY is not set")
	})
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		shouldPanic bool
	}{
		{"ClientID exists",      "GOOGLE_CLIENT_ID",     "dummy", false},
		{"ClientID missing",     "GOOGLE_CLIENT_ID",     "",      true},
		{"ClientSecret exists",  "GOOGLE_CLIENT_SECRET", "dummy", false},
		{"ClientSecret missing", "GOOGLE_CLIENT_SECRET", "",      true},
		{"RedirectURI exists",   "GOOGLE_REDIRECT_URI",  "dummy", false},
		{"RedirectURI missing",  "GOOGLE_REDIRECT_URI",  "",      true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.key, tt.value)

			if tt.shouldPanic {
				require.PanicsWithValue(t,
					fmt.Sprintf("missing environment variable: %s", tt.key),
					func() {
						config.GetEnv(tt.key)
					})
			} else {
				require.Equal(t, tt.value, config.GetEnv(tt.key))
			}
		})
	}
}
