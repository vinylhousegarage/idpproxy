package idpproxy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/logout"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/test/routertest"
)

type fakeLogoutSessionInvalidator struct {
	called    bool
	sessionID string
	err       error
}

func (f *fakeLogoutSessionInvalidator) Invalidate(
	_ context.Context,
	sessionID string,
) (*authsession.Session, error) {
	f.called = true
	f.sessionID = sessionID

	if f.err != nil {
		return nil, f.err
	}

	return &authsession.Session{}, nil
}

func TestLogoutRoute(t *testing.T) {
	t.Run("invalidates session and clears session cookie", func(t *testing.T) {
		logger := zap.NewNop()
		invalidator := &fakeLogoutSessionInvalidator{}

		d := routertest.NewMockRouterDeps(logger)
		d.Logout = logout.NewHandler(invalidator)

		r := router.NewRouter(d)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  authsession.SessionCookieName,
			Value: "session-123",
			Path:  "/",
		})
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNoContent, recorder.Code)
		require.True(t, invalidator.called)
		require.Equal(t, "session-123", invalidator.sessionID)

		cookie := findLogoutCookie(
			t,
			recorder.Result().Cookies(),
			authsession.SessionCookieName,
		)
		require.Equal(t, "/", cookie.Path)
		require.Equal(t, -1, cookie.MaxAge)
	})

	t.Run("returns no content when session cookie is missing", func(t *testing.T) {
		logger := zap.NewNop()
		invalidator := &fakeLogoutSessionInvalidator{}

		d := routertest.NewMockRouterDeps(logger)
		d.Logout = logout.NewHandler(invalidator)

		r := router.NewRouter(d)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNoContent, recorder.Code)
		require.False(t, invalidator.called)
	})
}

func findLogoutCookie(
	t *testing.T,
	cookies []*http.Cookie,
	name string,
) *http.Cookie {
	t.Helper()

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	t.Fatalf("cookie %q was not set", name)

	return nil
}
