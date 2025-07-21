package login_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestGoogleLoginRoute_RedirectsToGoogle(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	di := testhelpers.NewMockDeps(logger)

	r := router.NewRouter(di, http.FS(public.PublicFS))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/google/login", nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusFound, w.Code)
	require.Contains(t, w.Header().Get("Location"), "https://accounts.google.com/o/oauth2/v2/auth")
	require.Contains(t, w.Header().Get("Location"), "client_id=")
	require.Contains(t, w.Header().Get("Set-Cookie"), "oauth_state=")
}
