package me

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockVerifier struct {
	VerifyIDTokenFunc func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *mockVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return m.VerifyIDTokenFunc(ctx, idToken)
}

func TestMeHandler(t *testing.T) {
	t.Parallel()

	t.Run("OPTIONS preflight request", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodOptions, "/me", nil)
		w := httptest.NewRecorder()

		handler := NewMeHandler(nil, zap.NewNop())
		handler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		require.Equal(t, "GET, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
		require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
		require.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Authorization")
	})

	t.Run("GET with valid id_token", func(t *testing.T) {
		t.Parallel()

		mock := &mockVerifier{
			VerifyIDTokenFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{
					UID: "test-uid",
					Claims: map[string]interface{}{
						"iss": "https://issuer.example.com",
						"aud": "test-audience",
						"exp": float64(1234567890),
					},
				}, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer valid.token.here")
		w := httptest.NewRecorder()

		handler := NewMeHandler(mock, zap.NewNop())
		handler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "GET, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		require.JSONEq(t, `{
			"sub": "test-uid",
			"iss": "https://issuer.example.com",
			"aud": "test-audience",
			"exp": 1234567890
		}`, w.Body.String())
	})

	t.Run("GET with missing Authorization header", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		w := httptest.NewRecorder()

		handler := NewMeHandler(nil, zap.NewNop())
		handler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with invalid token format", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "invalid-format")
		w := httptest.NewRecorder()

		handler := NewMeHandler(nil, zap.NewNop())
		handler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, "GET, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with empty token after Bearer", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer ")
		w := httptest.NewRecorder()

		handler := NewMeHandler(nil, zap.NewNop())
		handler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Contains(t, w.Body.String(), `"error":`)
	})
}
