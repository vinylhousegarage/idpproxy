package config

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadIDTokenConfig(t *testing.T) {
	setValidIDTokenEnv := func(t *testing.T) {
		t.Helper()

		t.Setenv(
			"ID_TOKEN_ISSUER",
			"https://idpproxy.example.com",
		)
		t.Setenv(
			"ID_TOKEN_HMAC_KEY_BASE64",
			base64.StdEncoding.EncodeToString(
				[]byte("test-signing-key"),
			),
		)
		t.Setenv("ID_TOKEN_HMAC_KEY_ID", "test-hs256-key")
		t.Setenv("ID_TOKEN_TTL", "1h")
	}

	t.Run("loads ID token configuration", func(t *testing.T) {
		setValidIDTokenEnv(t)

		cfg, err := LoadIDTokenConfig()

		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Equal(
			t,
			"https://idpproxy.example.com",
			cfg.Issuer,
		)
		require.Equal(t, []byte("test-signing-key"), cfg.HMACKey)
		require.Equal(t, "test-hs256-key", cfg.HMACKeyID)
		require.Equal(t, time.Hour, cfg.TTL)
	})

	t.Run("trims surrounding spaces", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv(
			"ID_TOKEN_ISSUER",
			"  https://idpproxy.example.com  ",
		)
		t.Setenv(
			"ID_TOKEN_HMAC_KEY_ID",
			"  test-hs256-key  ",
		)
		t.Setenv("ID_TOKEN_TTL", "  30m  ")

		cfg, err := LoadIDTokenConfig()

		require.NoError(t, err)
		require.Equal(
			t,
			"https://idpproxy.example.com",
			cfg.Issuer,
		)
		require.Equal(t, "test-hs256-key", cfg.HMACKeyID)
		require.Equal(t, 30*time.Minute, cfg.TTL)
	})

	t.Run("returns error when issuer is not set", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_ISSUER", "")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.EqualError(t, err, "ID_TOKEN_ISSUER is not set")
	})

	t.Run("returns error when HMAC key is not set", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_HMAC_KEY_BASE64", "")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.EqualError(
			t,
			err,
			"ID_TOKEN_HMAC_KEY_BASE64 is not set",
		)
	})

	t.Run("returns error when HMAC key is invalid base64", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_HMAC_KEY_BASE64", "not-base64!!!")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.ErrorContains(
			t,
			err,
			"failed to decode ID_TOKEN_HMAC_KEY_BASE64",
		)
	})

	t.Run("returns error when HMAC key ID is not set", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_HMAC_KEY_ID", "")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.EqualError(
			t,
			err,
			"ID_TOKEN_HMAC_KEY_ID is not set",
		)
	})

	t.Run("returns error when TTL is not set", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_TTL", "")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.EqualError(t, err, "ID_TOKEN_TTL is not set")
	})

	t.Run("returns error when TTL is invalid", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_TTL", "invalid")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.ErrorContains(
			t,
			err,
			"ID_TOKEN_TTL is invalid",
		)
	})

	t.Run("returns error when TTL is not positive", func(t *testing.T) {
		setValidIDTokenEnv(t)
		t.Setenv("ID_TOKEN_TTL", "0s")

		cfg, err := LoadIDTokenConfig()

		require.Nil(t, cfg)
		require.EqualError(
			t,
			err,
			"ID_TOKEN_TTL must be greater than zero",
		)
	})
}
