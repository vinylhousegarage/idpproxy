package token

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

type fakeHTTPDoer struct {
	resp *http.Response
	err  error
}

func (f *fakeHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func newResp(status int, contentType, body string) *http.Response {
	rr := httptest.NewRecorder()
	rr.WriteHeader(status)
	if contentType != "" {
		rr.Header().Set("Content-Type", contentType)
	}
	rr.Body.WriteString(body)
	return rr.Result()
}

func TestExchange(t *testing.T) {
	t.Parallel()

	cfg := &config.GitHubOAuthConfig{
		ClientID:     "cid",
		ClientSecret: "sec",
		RedirectURI:  "http://localhost/callback",
	}

	ctx := context.Background()

	t.Run("success_returns_access_token", func(t *testing.T) {
		t.Parallel()

		body := `{"access_token":"gho_test","token_type":"bearer","scope":"read:user"}`
		httpc := &fakeHTTPDoer{
			resp: newResp(200, "application/json", body),
		}

		token, err := Exchange(ctx, httpc, cfg, "code123", "state123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "gho_test" {
			t.Fatalf("unexpected token: %s", token)
		}
	})

	t.Run("returns_error_when_http_do_fails", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPDoer{
			err: errors.New("network error"),
		}

		_, err := Exchange(ctx, httpc, cfg, "code123", "state123")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("returns_error_when_non_2xx_response", func(t *testing.T) {
		t.Parallel()

		httpc := &fakeHTTPDoer{
			resp: newResp(401, "text/plain", "unauthorized"),
		}

		_, err := Exchange(ctx, httpc, cfg, "code123", "state123")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("returns_error_when_oauth_error_response", func(t *testing.T) {
		t.Parallel()

		body := `{"error":"bad_verification_code","error_description":"expired"}`
		httpc := &fakeHTTPDoer{
			resp: newResp(200, "application/json", body),
		}

		_, err := Exchange(ctx, httpc, cfg, "code123", "state123")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("returns_error_when_access_token_missing", func(t *testing.T) {
		t.Parallel()

		body := `{"token_type":"bearer"}`
		httpc := &fakeHTTPDoer{
			resp: newResp(200, "application/json", body),
		}

		_, err := Exchange(ctx, httpc, cfg, "code123", "state123")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("uses_form_encoded_response", func(t *testing.T) {
		t.Parallel()

		body := "access_token=gho_form&scope=read%3Auser&token_type=bearer"
		httpc := &fakeHTTPDoer{
			resp: newResp(200, "application/x-www-form-urlencoded", body),
		}

		token, err := Exchange(ctx, httpc, cfg, "code123", "state123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "gho_form" {
			t.Fatalf("unexpected token: %s", token)
		}
	})
}
