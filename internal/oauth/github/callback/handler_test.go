package callback

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

func TestGitHubCallbackHandler_Serve(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tokenJSON := loadTestDataJSON(t, "testdata/token_success.json")
	userJSON := loadTestDataJSON(t, "testdata/user_success.json")

	t.Run("successfully_exchanges_code_saves_token_and_redirects", func(t *testing.T) {
		t.Parallel()

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

		h := newHandlerForTest(t, httpc, us, pcs, tokenRepo)

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
			t.Fatalf("expected 302, got=%d body=%s", rr.Code, rr.Body.String())
		}

		loc := rr.Header().Get("Location")
		if loc == "" {
			t.Fatal("missing Location header")
		}

		u, err := url.Parse(loc)
		if err != nil {
			t.Fatalf("invalid Location: %v (%s)", err, loc)
		}

		if got := u.Query().Get("code"); got != "proxycode-123" {
			t.Fatalf("expected proxycode-123, got=%s (loc=%s)", got, loc)
		}

		if got := u.Query().Get("state"); got != "st-abc" {
			t.Fatalf("expected state=st-abc, got=%s (loc=%s)", got, loc)
		}

		if !tokenRepo.called {
			t.Fatal("GitHubTokenRepo.Upsert was not called")
		}

		if tokenRepo.rec == nil {
			t.Fatal("GitHubTokenRepo.Upsert record is nil")
		}

		if got := tokenRepo.rec.FirebaseUID; got != "user-internal-123" {
			t.Fatalf("FirebaseUID = %q, want %q", got, "user-internal-123")
		}

		if got := tokenRepo.rec.Provider; got != "github" {
			t.Fatalf("Provider = %q, want %q", got, "github")
		}

		if tokenRepo.rec.GitHubID == "" {
			t.Fatal("GitHubID must not be empty")
		}

		if tokenRepo.rec.Login == "" {
			t.Fatal("Login must not be empty")
		}

		if tokenRepo.rec.AccessToken == "" {
			t.Fatal("AccessToken must not be empty")
		}

		if !pcs.called {
			t.Fatal("ProxyCodeService.Issue was not called")
		}

		assertStateCookieDeleted(t, rr)
	})

	t.Run("returns_400_when_state_is_invalid_and_deletes_cookie", func(t *testing.T) {
		t.Parallel()

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

		h := newHandlerForTest(t, httpc, us, pcs, tokenRepo)

		rr, req := newCallbackRequest(
			t,
			"/oauth/github/callback",
			"code123",
			"st-abc",
		)
		setStateCookie(req, "st-wrong")

		_, r := gin.CreateTestContext(rr)
		r.Use(apierror.ErrorLogger(h.OAuth.Logger))
		r.GET("/oauth/github/callback", h.Serve)
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got=%d body=%s", rr.Code, rr.Body.String())
		}

		resp := decodeErrorResponse(t, rr)
		if resp.Error != apierror.ErrorCodeInvalidState {
			t.Fatalf(
				"expected error=%s, got=%s",
				apierror.ErrorCodeInvalidState,
				resp.Error,
			)
		}

		if tokenRepo.called {
			t.Fatal("GitHubTokenRepo.Upsert must not be called on invalid state")
		}

		if pcs.called {
			t.Fatal("ProxyCodeService.Issue must not be called on invalid state")
		}

		assertStateCookieDeleted(t, rr)
	})

	t.Run("returns_502_when_token_exchange_fails", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPClient{
			tokenJSON:     tokenJSON,
			userJSON:      userJSON,
			forceTokenErr: true,
		}
		us := &fakeUserService{
			returnID: "user-internal-123",
		}
		pcs := &fakeProxyCodeService{
			proxyCode: "proxycode-123",
		}
		tokenRepo := &fakeGitHubTokenRepo{}

		h := newHandlerForTest(t, httpc, us, pcs, tokenRepo)

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

		if rr.Code != http.StatusBadGateway {
			t.Fatalf("expected 502, got=%d body=%s", rr.Code, rr.Body.String())
		}

		resp := decodeErrorResponse(t, rr)
		if resp.Error != apierror.ErrorCodeGitHubTokenRequest {
			t.Fatalf(
				"expected error=%s, got=%s",
				apierror.ErrorCodeGitHubTokenRequest,
				resp.Error,
			)
		}

		if tokenRepo.called {
			t.Fatal("GitHubTokenRepo.Upsert must not be called when token exchange fails")
		}

		if pcs.called {
			t.Fatal("ProxyCodeService.Issue must not be called when token exchange fails")
		}
	})

	t.Run("returns_500_when_GitHub_token_save_fails", func(t *testing.T) {
		t.Parallel()

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
		tokenRepo := &fakeGitHubTokenRepo{
			err: errors.New("Firestore unavailable"),
		}

		h := newHandlerForTest(t, httpc, us, pcs, tokenRepo)

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
			t.Fatalf("expected 500, got=%d body=%s", rr.Code, rr.Body.String())
		}

		resp := decodeErrorResponse(t, rr)
		if resp.Error != apierror.ErrorCodeGitHubTokenUpsert {
			t.Fatalf(
				"expected error=%s, got=%s",
				apierror.ErrorCodeGitHubTokenUpsert,
				resp.Error,
			)
		}

		if !tokenRepo.called {
			t.Fatal("GitHubTokenRepo.Upsert was not called")
		}

		if pcs.called {
			t.Fatal("ProxyCodeService.Issue must not be called when token save fails")
		}

		assertStateCookieDeleted(t, rr)
	})
}
