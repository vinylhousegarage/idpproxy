package callback

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildTokenRequestBody_Success(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	t.Setenv("GOOGLE_REDIRECT_URI", "https://localhost/callback")

	code := "auth-code-123"

	body, err := BuildTokenRequestBody(code)
	require.NoError(t, err)

	values, err := url.ParseQuery(body)
	require.NoError(t, err)
	require.Equal(t, "authorization_code", values.Get("grant_type"))
	require.Equal(t, code, values.Get("code"))
	require.Equal(t, "test-client-id", values.Get("client_id"))
	require.Equal(t, "test-client-secret", values.Get("client_secret"))
	require.Equal(t, "https://localhost/callback", values.Get("redirect_uri"))
}

func TestBuildTokenRequestBody_MissingEnv(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	t.Setenv("GOOGLE_REDIRECT_URI", "")

	body, err := BuildTokenRequestBody("dummy")
	require.Error(t, err)
	require.Empty(t, body)
	require.Contains(t, err.Error(), "missing required environment variables")
}
