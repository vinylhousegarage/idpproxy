package idpproxy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	githubcallback "github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	idpproxytoken "github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/token"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

type whoAmISessionValidator struct {
	session *authsession.Session
	err     error

	called    bool
	sessionID string
}

func (v *whoAmISessionValidator) Validate(
	_ context.Context,
	sessionID string,
) (*authsession.Session, error) {
	v.called = true
	v.sessionID = sessionID

	if v.err != nil {
		return nil, v.err
	}

	return v.session, nil
}

func TestWhoAmIRoute(t *testing.T) {
	t.Parallel()

	t.Run("returns_401_without_session_cookie", func(t *testing.T) {
		t.Parallel()

		validator := &whoAmISessionValidator{}
		r := newWhoAmIRouter(t, validator)

		req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		require.False(t, validator.called)
	})

	t.Run("returns_authenticated_user_id", func(t *testing.T) {
		t.Parallel()

		validator := &whoAmISessionValidator{
			session: &authsession.Session{
				SessionID: "session-123",
				UserID:    "firebase-user-123",
				ExpiresAt: time.Now().Add(time.Hour),
			},
		}
		r := newWhoAmIRouter(t, validator)

		req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
		req.AddCookie(&http.Cookie{
			Name:  authsession.SessionCookieName,
			Value: "session-123",
			Path:  "/",
		})
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.True(t, validator.called)
		require.Equal(t, "session-123", validator.sessionID)
		require.JSONEq(
			t,
			`{"user_id":"firebase-user-123"}`,
			recorder.Body.String(),
		)
	})
}

func newWhoAmIRouter(
	t *testing.T,
	sessionValidator authsession.SessionValidator,
) http.Handler {
	t.Helper()

	logger := zap.NewNop()

	d := router.NewRouterDeps(
		public.PublicFS,
		testhelpers.NewMockGitHubAPIDeps(logger),
		testhelpers.NewMockGitHubOAuthDeps(logger),
		&githubcallback.GitHubCallbackHandler{},
		&githubme.GitHubMeHandler{},
		sessionValidator,
		idpproxytoken.NewHandler(
			&idpproxytoken.Service{},
			logger,
		),
		testhelpers.NewMockGoogleDeps(logger),
		logger,
		testhelpers.NewMockSystemDeps(logger),
	)

	return router.NewRouter(d)
}
