package callback

import (
	"net/http"
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
	userJSON  := loadTestDataJSON(t, "testdata/user_success.json")

	httpc := &fakeHTTPClient{tokenJSON: tokenJSON, userJSON: userJSON}
	us    := &fakeUserService{returnID: "user-internal-123"}
	iss   := &fakeIssuer{jwt: "jwt.mock", kid: "kid-1"}

	h := newHandlerForTest(t, httpc, us, iss)

	rr, req := newCallbackRequest(t, "/oauth/github/callback", "code123", "st-abc")
	setStateCookie(req, "st-abc")

	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	h.Serve(ctx)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got=%d body=%s", rr.Code, rr.Body.String())
	}

	if got := rr.Body.String(); !containsAll(got, `"ok":true`, `"id_token":"jwt.mock"`, `"provider":"github"`) {
		t.Fatalf("unexpected body: %s", got)
	}
}
