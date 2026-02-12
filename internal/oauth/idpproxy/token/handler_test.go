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

		handler := NewHandler(svc, zap.NewNop())

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

		if resp.IDToken == "" {
			t.Fatal("id_token should not be empty")
		}
	})

	t.Run("invalid grant returns oauth error", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()
		handler := NewHandler(svc, zap.NewNop())

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
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

		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString("{bad json"))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
