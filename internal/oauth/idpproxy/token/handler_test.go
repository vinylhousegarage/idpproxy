package token

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestTokenHandler(t *testing.T) {
	t.Parallel()

	validReq := TokenRequest{
		GrantType: "authorization_code",
		Code:      "valid",
		ClientID:  "client-1",
	}

	t.Run("success returns 200 and id_token response", func(t *testing.T) {
		t.Parallel()

		now := time.Now()

		svc := &Service{
			Store: &mockStore{
				code: &AuthCode{
					UserID:    "user1",
					ClientID:  "client-1",
					ExpiresAt: now.Add(time.Hour),
				},
			},
			Clock: fixedClock{t: now},
			IDTokenIssuer: &mockIDTokenIssuer{
				token: "signed-id-token",
			},
			IDTokenTTL: time.Hour,
		}

		handler := NewHandler(svc, zap.NewNop())

		body, err := json.Marshal(validReq)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/token",
			bytes.NewReader(body),
		)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var resp TokenResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode error: %v", err)
		}

		if resp.IDToken != "signed-id-token" {
			t.Fatalf(
				"id_token = %q, want %q",
				resp.IDToken,
				"signed-id-token",
			)
		}
	})

	t.Run("invalid grant returns oauth error", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()
		handler := NewHandler(svc, zap.NewNop())

		body, err := json.Marshal(validReq)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/token",
			bytes.NewReader(body),
		)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}

		var oauthErr map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&oauthErr); err != nil {
			t.Fatalf("decode error: %v", err)
		}

		if oauthErr["error"] != "invalid_grant" {
			t.Fatalf("unexpected error: %v", oauthErr)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()
		handler := NewHandler(svc, zap.NewNop())

		req := httptest.NewRequest(
			http.MethodPost,
			"/token",
			bytes.NewBufferString("{bad json"),
		)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
