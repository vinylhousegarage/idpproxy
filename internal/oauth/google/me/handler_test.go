package me

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/apperror"
)

type mockVerifier struct {
	VerifyIDTokenFunc func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *mockVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return m.VerifyIDTokenFunc(ctx, idToken)
}

type brokenWriter struct {
	gin.ResponseWriter
}

func (bw *brokenWriter) Write(p []byte) (int, error) {
	return 0, apperror.New(http.StatusInternalServerError, "write failed")
}

func TestMeHandler(t *testing.T) {
	t.Parallel()

	t.Run("OPTIONS preflight request", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodOptions, "/me", nil)

		handler := NewMeHandler(nil, zap.NewNop())
		handler.Serve(c)

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

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
		c.Request.Header.Set("Authorization", "Bearer valid.token.here")

		handler := NewMeHandler(mock, zap.NewNop())
		handler.Serve(c)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
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

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)

		handler := NewMeHandler(nil, zap.NewNop())
		handler.Serve(c)

		require.Equal(t, ErrMissingAuthorizationHeader.Code, w.Code)
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with invalid token format", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
		c.Request.Header.Set("Authorization", "invalid-format")

		handler := NewMeHandler(nil, zap.NewNop())
		handler.Serve(c)

		require.Equal(t, ErrInvalidAuthorizationHeaderFormat.Code, w.Code)
		require.Equal(t, "GET, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with empty token after Bearer", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
		c.Request.Header.Set("Authorization", "Bearer ")

		handler := NewMeHandler(nil, zap.NewNop())
		handler.Serve(c)

		require.Equal(t, ErrEmptyBearerToken.Code, w.Code)
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with JSON encoding error", func(t *testing.T) {
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

		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		bw := &brokenWriter{ResponseWriter: c.Writer}
		c.Writer = bw
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
		c.Request.Header.Set("Authorization", "Bearer valid.token.here")

		handler := NewMeHandler(mock, zap.NewNop())
		handler.Serve(c)

		require.Equal(t, ErrFailedToWriteUserResponse.Code, rr.Code)
		require.Equal(t, "", rr.Body.String())
	})
}
