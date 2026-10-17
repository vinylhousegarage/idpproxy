package github_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	githubcallback "github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

type fakeGitHubMeService struct {
	user *response.GitHubUserAPIResponse
	err  error

	called      bool
	firebaseUID string
}

func (s *fakeGitHubMeService) Get(
	_ context.Context,
	firebaseUID string,
) (*response.GitHubUserAPIResponse, error) {
	s.called = true
	s.firebaseUID = firebaseUID

	if s.err != nil {
		return nil, s.err
	}

	return s.user, nil
}

type fakeSessionValidator struct {
	session *authsession.Session
	err     error

	called    bool
	sessionID string
}

func (v *fakeSessionValidator) Validate(
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

func TestGitHubMeRoute(t *testing.T) {
	t.Parallel()

	t.Run("returns_401_without_session_cookie", func(t *testing.T) {
		t.Parallel()

		service := &fakeGitHubMeService{
			user: &response.GitHubUserAPIResponse{},
		}
		validator := &fakeSessionValidator{
			session: &authsession.Session{
				UserID: "firebase-user-123",
			},
		}
		r := newGitHubMeRouter(t, service, validator)

		req := httptest.NewRequest(http.MethodGet, "/github/me", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.False(t, validator.called)
		require.False(t, service.called)
	})

	t.Run("returns_200_and_calls_service_with_session_user_ID", func(t *testing.T) {
		t.Parallel()

		service := &fakeGitHubMeService{
			user: &response.GitHubUserAPIResponse{
				Login: "octocat",
			},
		}
		validator := &fakeSessionValidator{
			session: &authsession.Session{
				SessionID: "session-123",
				UserID:    "firebase-user-123",
			},
		}
		r := newGitHubMeRouter(t, service, validator)

		req := httptest.NewRequest(http.MethodGet, "/github/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  authsession.SessionCookieName,
			Value: "session-123",
			Path:  "/",
		})
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, validator.called)
		require.Equal(t, "session-123", validator.sessionID)
		require.True(t, service.called)
		require.Equal(t, "firebase-user-123", service.firebaseUID)
		require.Contains(t, rr.Body.String(), `"login":"octocat"`)
	})
}

func newGitHubMeRouter(
	t *testing.T,
	service *fakeGitHubMeService,
	validator authsession.SessionValidator,
) *gin.Engine {
	t.Helper()

	logger := zap.NewNop()

	d := router.NewRouterDeps(
		public.PublicFS,
		testhelpers.NewMockGitHubAPIDeps(logger),
		testhelpers.NewMockGitHubOAuthDeps(logger),
		&githubcallback.GitHubCallbackHandler{},
		githubme.NewGitHubMeHandler(service),
		validator,
		testhelpers.NewMockGoogleDeps(logger),
		logger,
		testhelpers.NewMockSystemDeps(logger),
	)

	return router.NewRouter(d)
}
