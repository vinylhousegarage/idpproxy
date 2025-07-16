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
