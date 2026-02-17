package callback

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func newHandlerForTest(
	t *testing.T,
	httpc *fakeHTTPClient,
	us *fakeUserService,
	iss *fakeIssuer,
	acs *fakeAuthCodeService,
) *GitHubCallbackHandler {
	t.Helper()

	logger := zaptest.NewLogger(t)

	oauth := deps.NewGitHubOAuthDeps(
		&config.GitHubOAuthConfig{
			ClientID:     "cid",
			ClientSecret: "sec",
			RedirectURI:  "http://localhost/cb",
		},
		logger,
	)

	api := deps.NewGitHubAPIDeps(
		&config.GitHubAPIConfig{
			APIVersion: config.GitHubAPIVersion,
			BaseURL:    "https://api.github.com",
			UserAgent:  config.UserAgent(),
		},
		httpc,
		logger,
	)

	return NewGitHubCallbackHandler(
		oauth,
		api,
		us,
		&issuerAdapter{iss},
		acs, // ← 追加
		"test-client",
	)
}

func newCallbackRequest(t *testing.T, path, code, state string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	q := url.Values{}
	if code != "" {
		q.Set("code", code)
	}
	if state != "" {
		q.Set("state", state)
	}

	req := httptest.NewRequest(http.MethodGet, path+"?"+q.Encode(), nil)
	req = req.WithContext(context.Background())

	return httptest.NewRecorder(), req
}

func setStateCookie(r *http.Request, state string) {
	r.AddCookie(&http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60,
	})
}

func loadTestDataJSON(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func assertStateCookieDeleted(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()

	setCookies := rr.Header().Values("Set-Cookie")
	if len(setCookies) == 0 {
		t.Fatalf("expected Set-Cookie header, got none")
	}

	joined := strings.Join(setCookies, "\n")

	if !strings.Contains(joined, stateCookieName) {
		t.Fatalf("expected state cookie deletion, got: %s", joined)
	}

	if !(strings.Contains(joined, "Max-Age=0") ||
		strings.Contains(strings.ToLower(joined), "expires=")) {
		t.Fatalf("expected deletion attributes, got: %s", joined)
	}
}
