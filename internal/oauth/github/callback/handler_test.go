package callback

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGitHubCallbackHandler_Serve(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tokenJSON := loadTestDataJSON(t, "testdata/token_success.json")
	userJSON := loadTestDataJSON(t, "testdata/user_success.json")

	t.Run("successfully_exchanges_code_and_redirects_and_deletes_state_cookie", func(t *testing.T) {
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
			t.Fatalf("expected 302, got=%d body=%s", rr.Code, rr.Body.String())
		}

		loc := rr.Header().Get("Location")
		if loc == "" {
			t.Fatalf("missing Location header")
		}

		u, err := url.Parse(loc)
		if err != nil {
			t.Fatalf("invalid Location: %v (%s)", err, loc)
		}

		if u.Query().Get("code") == "" {
			t.Fatalf("missing code param in redirect: %s", loc)
		}
		if u.Query().Get("state") == "" {
			t.Fatalf("missing state param in redirect: %s", loc)
		}

		assertStateCookieDeleted(t, rr)
	})

	t.Run("returns_400_when_state_is_invalid_and_deletes_cookie", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPClient{tokenJSON: tokenJSON, userJSON: userJSON}
		us := &fakeUserService{returnID: "user-internal-123"}
		iss := &fakeIssuer{jwt: "jwt.mock", kid: "kid-1"}
		h := newHandlerForTest(t, httpc, us, iss)

		rr, req := newCallbackRequest(t, "/oauth/github/callback", "code123", "st-abc")
		setStateCookie(req, "st-wrong")

		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = req

		h.Serve(ctx)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got=%d body=%s", rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), `"error":"invalid state"`) {
			t.Fatalf("unexpected body: %s", rr.Body.String())
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
		us := &fakeUserService{returnID: "user-internal-123"}
		iss := &fakeIssuer{jwt: "jwt.mock", kid: "kid-1"}
		h := newHandlerForTest(t, httpc, us, iss)

		rr, req := newCallbackRequest(t, "/oauth/github/callback", "code123", "st-abc")
		setStateCookie(req, "st-abc")

		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = req

		h.Serve(ctx)

		if rr.Code != http.StatusBadGateway {
			t.Fatalf("expected 502, got=%d body=%s", rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), `"error":"token request failed"`) {
			t.Fatalf("unexpected body: %s", rr.Body.String())
		}
	})
}
