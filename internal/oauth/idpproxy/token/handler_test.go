package token

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type stubService struct {
	resp *TokenResponse
	err  error
}

func (s *stubService) Exchange(ctx context.Context, req TokenRequest) (*TokenResponse, error) {
	return s.resp, s.err
}

func TestTokenHandler(t *testing.T) {
	t.Parallel()

	validReq := TokenRequest{
		GrantType:    "authorization_code",
		Code:         "valid",
		ClientID:     "client-1",
		ClientSecret: "secret",
	}

	t.Run("success returns 200 and token response", func(t *testing.T) {
		t.Parallel()

		svc := &Service{
			Store: &mockStore{
				code: &AuthCode{
					UserID:    "user1",
					ClientID:  "client-1",
					ExpiresAt: time.Now().Add(time.Hour),
				},
			},
			Clock: fixedClock{t: time.Now()},
		}

		handler := TokenHandler(svc)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var resp TokenResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode error: %v", err)
		}

		if resp.AccessToken == "" {
			t.Fatal("access_token should not be empty")
		}
	})

	t.Run("invalid grant returns oauth error", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()

		handler := TokenHandler(svc)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}

		var oauthErr map[string]string
		json.NewDecoder(rec.Body).Decode(&oauthErr)

		if oauthErr["error"] != "invalid_grant" {
			t.Fatalf("unexpected error: %v", oauthErr)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()
		handler := TokenHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString("{bad json"))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
