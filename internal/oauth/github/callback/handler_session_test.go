package callback

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

func TestGitHubCallbackHandler_Serve_Session(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenJSON := loadTestDataJSON(t, "testdata/token_success.json")
	userJSON := loadTestDataJSON(t, "testdata/user_success.json")

	t.Run("starts_session_sets_cookie_and_issues_proxy_code", func(t *testing.T) {
		expiresAt := time.Date(
			2026,
			time.October,
			6,
			0,
			0,
			0,
			0,
			time.UTC,
		)

		httpc := &fakeHTTPClient{
			tokenJSON: tokenJSON,
			userJSON:  userJSON,
		}
		us := &fakeUserService{
			returnID: "user-internal-123",
		}
		pcs := &fakeProxyCodeService{
			proxyCode: "proxycode-123",
		}
		tokenRepo := &fakeGitHubTokenRepo{}
		sessionSvc := &fakeSessionService{
			session: &authsession.Session{
				SessionID: "session-123",
				UserID:    "user-internal-123",
				ExpiresAt: expiresAt,
			},
		}

		h := newHandlerForTest(
			t,
			httpc,
			us,
			pcs,
			tokenRepo,
			sessionSvc,
		)

		rr, req := newCallbackRequest(
			t,
			"/oauth/github/callback",
			"code123",
			"st-abc",
		)
		setStateCookie(req, "st-abc")

		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = req

		h.Serve(ctx)

		if rr.Code != http.StatusFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusFound)
		}

		if !tokenRepo.called {
			t.Fatal("GitHubTokenRepo.Upsert was not called")
		}

		if !sessionSvc.called {
			t.Fatal("SessionService.Start was not called")
		}

		if got, want := sessionSvc.userID, "user-internal-123"; got != want {
			t.Fatalf("SessionService.Start userID = %q, want %q", got, want)
		}

		assertSessionCookieSet(
			t,
			rr,
			"session-123",
			expiresAt,
		)

		if !pcs.called {
			t.Fatal("ProxyCodeService.Issue was not called")
		}

		assertStateCookieDeleted(t, rr)
	})

	t.Run("returns_500_when_session_start_fails", func(t *testing.T) {
		httpc := &fakeHTTPClient{
			tokenJSON: tokenJSON,
			userJSON:  userJSON,
		}
		us := &fakeUserService{
			returnID: "user-internal-123",
		}
		pcs := &fakeProxyCodeService{
			proxyCode: "proxycode-123",
		}
		tokenRepo := &fakeGitHubTokenRepo{}
		sessionSvc := &fakeSessionService{
			err: errors.New("session store unavailable"),
		}

		h := newHandlerForTest(
			t,
			httpc,
			us,
			pcs,
			tokenRepo,
			sessionSvc,
		)

		rr, req := newCallbackRequest(
			t,
			"/oauth/github/callback",
			"code123",
			"st-abc",
		)
		setStateCookie(req, "st-abc")

		_, r := gin.CreateTestContext(rr)
		r.Use(apierror.ErrorLogger(h.OAuth.Logger))
		r.GET("/oauth/github/callback", h.Serve)
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Fatalf(
				"status = %d, want %d",
				rr.Code,
				http.StatusInternalServerError,
			)
		}

		resp := decodeErrorResponse(t, rr)
		if resp.Error != apierror.ErrorCodeSessionStart {
			t.Fatalf(
				"error = %q, want %q",
				resp.Error,
				apierror.ErrorCodeSessionStart,
			)
		}

		if !tokenRepo.called {
			t.Fatal("GitHubTokenRepo.Upsert was not called")
		}

		if !sessionSvc.called {
			t.Fatal("SessionService.Start was not called")
		}

		if pcs.called {
			t.Fatal(
				"ProxyCodeService.Issue must not be called when session start fails",
			)
		}

		assertStateCookieDeleted(t, rr)
	})
}

func assertSessionCookieSet(
	t *testing.T,
	rr *httptest.ResponseRecorder,
	sessionID string,
	expiresAt time.Time,
) {
	t.Helper()

	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name != authsession.SessionCookieName {
			continue
		}

		if got, want := cookie.Value, sessionID; got != want {
			t.Fatalf("session cookie value = %q, want %q", got, want)
		}

		if got, want := cookie.Path, "/"; got != want {
			t.Fatalf("session cookie path = %q, want %q", got, want)
		}

		if !cookie.HttpOnly {
			t.Fatal("session cookie must be HttpOnly")
		}

		if !cookie.Secure {
			t.Fatal("session cookie must be Secure")
		}

		if !cookie.Expires.Equal(expiresAt) {
			t.Fatalf(
				"session cookie expires = %v, want %v",
				cookie.Expires,
				expiresAt,
			)
		}

		return
	}

	t.Fatalf(
		"session cookie %q was not set",
		authsession.SessionCookieName,
	)
}
