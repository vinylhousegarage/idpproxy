package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetIDTokenCookie(t *testing.T) {
	t.Parallel()

	idToken := "dummy.token.value"
	rr := httptest.NewRecorder()

	SetIDTokenCookie(rr, idToken)

	cookies := rr.Result().Cookies()
	require.Len(t, cookies, 1)

	cookie := cookies[0]
	require.Equal(t, "id_token", cookie.Name)
	require.Equal(t, idToken, cookie.Value)
	require.True(t, cookie.HttpOnly)
	require.True(t, cookie.Secure)
	require.Equal(t, "/", cookie.Path)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)

	now := time.Now()
	require.WithinDuration(t, now.Add(15*time.Minute), cookie.Expires, time.Minute)
}
