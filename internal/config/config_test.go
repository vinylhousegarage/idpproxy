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
