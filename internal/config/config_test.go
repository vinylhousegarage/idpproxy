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
