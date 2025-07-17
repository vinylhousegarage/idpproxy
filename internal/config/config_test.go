package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPort_WithEnvVar(t *testing.T) {
	os.Setenv("PORT", "12345")
	defer os.Unsetenv("PORT")

	require.Equal(t, "12345", GetPort(), "PORT should match expected value")
}

func TestGetPort_WithoutEnvVar(t *testing.T) {
	os.Unsetenv("PORT")

	require.Equal(t, "9000", GetPort(), "PORT should match default value")
}

func TestGetOpenAPIURL(t *testing.T) {
	t.Setenv("OPENAPI_URL", "https://test.example.com/openapi.json")
	require.Equal(t, "https://test.example.com/openapi.json", GetOpenAPIURL())
}

func TestLoadGoogleConfig_Success(t *testing.T) {
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:9000/callback")
	os.Setenv("GOOGLE_RESPONSE_TYPE", "code")
	os.Setenv("GOOGLE_SCOPE", "unused")
	os.Setenv("GOOGLE_ACCESS_TYPE", "unused")
	os.Setenv("GOOGLE_PROMPT", "unused")

	cfg, err := LoadGoogleConfig()
	require.NoError(t, err)
	require.Equal(t, "test-client-id", cfg.ClientID)
	require.Equal(t, "test-client-secret", cfg.ClientSecret)
	require.Equal(t, "http://localhost:9000/callback", cfg.RedirectURI)
	require.Equal(t, "code", cfg.ResponseType)
	require.Equal(t, "openid email profile", cfg.Scope)
	require.Equal(t, "offline", cfg.AccessType)
	require.Equal(t, "consent", cfg.Prompt)
}

func TestLoadGoogleConfig_MissingVariables(t *testing.T) {
	os.Clearenv()

	cfg, err := LoadGoogleConfig()
	require.Nil(t, cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required environment variables")
}
