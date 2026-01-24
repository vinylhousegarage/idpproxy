package token

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestExtractAccessTokenFromResponse(t *testing.T) {
	t.Parallel()

	newResp := func(status int, ct string, body string) *http.Response {
		t.Helper()

		h := make(http.Header)
		if ct != "" {
			h.Set("Content-Type", ct)
		}
		return &http.Response{
			StatusCode: status,
			Header:     h,
			Body:       io.NopCloser(strings.NewReader(body)),
		}
	}

	t.Run("JSON_Success", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/json; charset=utf-8",
			`{"access_token":"gho_abc","token_type":"bearer","scope":"read:user"}`)
		token, err := ExtractAccessTokenFromResponse(resp)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token != "gho_abc" {
			t.Fatalf("token mismatch: got %q", token)
		}
	})

	t.Run("JSON_Success_WithTrailingSpaces", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/json", `{"access_token":"  gho_trim  \n"}`)
		token, err := ExtractAccessTokenFromResponse(resp)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token != "gho_trim" {
			t.Fatalf("got %q", token)
		}
	})

	t.Run("JSON_ErrorFields", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/json",
			`{"error":"bad_verification_code","error_description":"expired"}`)
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrGitHubOAuthError) {
			t.Fatalf("expected ErrGitHubOAuthError, got %v", err)
		}
	})

	t.Run("JSON_MissingToken", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/json", `{}`)
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrMissingAccessToken) {
			t.Fatalf("expected ErrMissingAccessToken, got %v", err)
		}
	})

	t.Run("Form_Success", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/x-www-form-urlencoded",
			"access_token=gho_xyz&scope=read%3Auser&token_type=bearer")
		token, err := ExtractAccessTokenFromResponse(resp)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token != "gho_xyz" {
			t.Fatalf("token mismatch: got %q", token)
		}
	})

	t.Run("Form_ErrorFields", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/x-www-form-urlencoded",
			"error=incorrect_client_credentials&error_description=nope")
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrGitHubOAuthError) {
			t.Fatalf("expected ErrGitHubOAuthError, got %v", err)
		}
	})

	t.Run("Form_MissingToken", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/x-www-form-urlencoded",
			"scope=read%3Auser&token_type=bearer")
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrMissingAccessToken) {
			t.Fatalf("expected ErrMissingAccessToken, got %v", err)
		}
	})

	t.Run("NoContentType_TreatedAsForm", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "", "access_token=gho_nct&token_type=bearer")
		token, err := ExtractAccessTokenFromResponse(resp)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token != "gho_nct" {
			t.Fatalf("token mismatch: got %q", token)
		}
	})

	t.Run("EmptyBody", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/json", "")
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrEmptyBody) {
			t.Fatalf("expected ErrEmptyBody, got %v", err)
		}
	})

	t.Run("Non2xx_ShowsSnippet", func(t *testing.T) {
		t.Parallel()

		long := bytes.Repeat([]byte("X"), 5000)
		resp := newResp(401, "text/plain", string(long))
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrNon2xxStatus) {
			t.Fatalf("expected ErrNon2xxStatus, got %v", err)
		}
	})

	t.Run("ContentTypeWithWeirdCasing_JSONDetected", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "Application/JSON; Charset=UTF-8", `{"access_token":"gho_case"}`)
		token, err := ExtractAccessTokenFromResponse(resp)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token != "gho_case" {
			t.Fatalf("token mismatch: got %q", token)
		}
	})

	t.Run("Form_InvalidEncoding_ParseError", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/x-www-form-urlencoded", "access_token=gho&scope=%ZZ")
		_, err := ExtractAccessTokenFromResponse(resp)
		if !errors.Is(err, ErrParseFormBody) {
			t.Fatalf("expected ErrParseFormBody, got %v", err)
		}
	})

	t.Run("VendorJSON_TreatedAsJSON", func(t *testing.T) {
		t.Parallel()

		resp := newResp(200, "application/vnd.github+json", `{"access_token":"gho_vjson"}`)
		token, err := ExtractAccessTokenFromResponse(resp)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if token != "gho_vjson" {
			t.Fatalf("token mismatch: got %q", token)
		}
	})
}
