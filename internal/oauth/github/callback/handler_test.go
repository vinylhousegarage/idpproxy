package callback

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

func TestGitHubCallbackHandler_Serve(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tokenJSON := loadTestDataJSON(t, "testdata/token_success.json")
	userJSON := loadTestDataJSON(t, "testdata/user_success.json")

	t.Run("success_redirects_with_code_and_state_and_deletes_cookie", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPClient{tokenJSON: tokenJSON, userJSON: userJSON}
		us := &fakeUserService{returnID: "user-internal-123"}
		iss := &fakeIssuer{jwt: "jwt.mock", kid: "kid-1"}

		h := newHandlerForTest(t, httpc, us, iss)

		rr, req := newCallbackRequest(t, "/oauth/github/callback", "code123", "st-abc")
		setStateCookie(req, "st-abc")

		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = req

		h.Serve(ctx)

		if rr.Code != http.StatusFound {
			t.Fatalf("unexpected status: got=%d body=%s", rr.Code, rr.Body.String())
		}

		loc := rr.Header().Get("Location")
		if loc == "" {
			t.Fatalf("missing Location header")
		}
		if !strings.Contains(loc, "code=") {
			t.Fatalf("missing code in Location: %s", loc)
		}
		if !strings.Contains(loc, "state=") {
			t.Fatalf("missing state in Location: %s", loc)
		}

		u, err := url.Parse(loc)
		if err != nil {
			t.Fatalf("invalid Location: %v (%s)", err, loc)
		}
		q := u.Query()
		if q.Get("code") == "" {
			t.Fatalf("missing code query param: %s", loc)
		}
		if q.Get("state") == "" {
			t.Fatalf("missing state query param: %s", loc)
		}

		setCookies := rr.Header().Values("Set-Cookie")
		if len(setCookies) == 0 {
			t.Fatalf("expected Set-Cookie header to delete state cookie, got none")
		}

		joined := strings.Join(setCookies, "\n")
		if !containsAll(joined, stateCookieName) {
			t.Fatalf("expected state cookie deletion, Set-Cookie=%s", joined)
		}
		if !(strings.Contains(joined, "Max-Age=0") || strings.Contains(strings.ToLower(joined), "expires=")) {
			t.Fatalf("expected cookie deletion attributes (Max-Age=0 or Expires=...), Set-Cookie=%s", joined)
		}
	})

	t.Run("state_mismatch_returns_400_invalid_state_and_deletes_cookie", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPClient{tokenJSON: tokenJSON, userJSON: userJSON}
		us := &fakeUserService{returnID: "user-internal-123"}
		iss := &fakeIssuer{jwt: "jwt.mock", kid: "kid-1"}

		h := newHandlerForTest(t, httpc, us, iss)

		rr, req := newCallbackRequest(t, "/oauth/github/callback", "code123", "st-abc")
		setStateCookie(req, "st-different")

		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = req

		h.Serve(ctx)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status: got=%d body=%s", rr.Code, rr.Body.String())
		}
		if got := rr.Body.String(); !strings.Contains(got, `"error":"invalid state"`) {
			t.Fatalf("unexpected body: %s", got)
		}

		setCookies := rr.Header().Values("Set-Cookie")
		if len(setCookies) == 0 {
			t.Fatalf("expected Set-Cookie header to delete state cookie, got none")
		}
		joined := strings.Join(setCookies, "\n")
		if !containsAll(joined, stateCookieName) {
			t.Fatalf("expected state cookie deletion, Set-Cookie=%s", joined)
		}
	})

	t.Run("token_exchange_http_error_returns_502", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPClient{
			tokenJSON:     tokenJSON,
			userJSON:      userJSON,
			forceTokenErr: true,
		}
		us := &fakeUserService{returnID: "user-internal-123"}
		iss := &fakeIssuer{jwt: "jwt.mock", kid: "kid-1"}

		h := newHandlerForTest(t, httpc, us, iss)

		rr, req := newCallbackRequest(t, "/oauth/github/callback", "code123", "st-abc")
		setStateCookie(req, "st-abc")

		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = req

		h.Serve(ctx)

		if rr.Code != http.StatusBadGateway {
			t.Fatalf("unexpected status: got=%d body=%s", rr.Code, rr.Body.String())
		}
		if got := rr.Body.String(); !strings.Contains(got, `"error":"token request failed"`) {
			t.Fatalf("unexpected body: %s", got)
		}
	})
}
