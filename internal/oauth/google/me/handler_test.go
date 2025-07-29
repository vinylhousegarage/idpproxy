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
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/verify"
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

func setupTestRouter(verifier verify.Verifier, logger *zap.Logger) *gin.Engine {
	router := gin.New()
	handler := NewMeHandler(verifier, logger)
	router.GET("/me", handler.Serve)
	router.OPTIONS("/me", handler.Serve)
	return router
}

func newRequest(method, path, authHeader string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	return req, httptest.NewRecorder()
}

func TestMeHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("OPTIONS preflight request", func(t *testing.T) {
		t.Parallel()
		router := setupTestRouter(nil, zap.NewNop())

		req, w := newRequest(http.MethodOptions, "/me", "")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNoContent, w.Code)
		require.Equal(t, "GET, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		require.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
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

		router := setupTestRouter(mock, zap.NewNop())
		req, w := newRequest(http.MethodGet, "/me", "Bearer valid.token.here")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		require.JSONEq(t, `{
			"sub": "test-uid",
			"iss": "https://issuer.example.com",
			"aud": "test-audience",
			"exp": 1234567890
		}`, w.Body.String())
	})

	t.Run("GET with missing Authorization header", func(t *testing.T) {
		t.Parallel()

		router := setupTestRouter(nil, zap.NewNop())
		req, w := newRequest(http.MethodGet, "/me", "")
		router.ServeHTTP(w, req)

		require.Equal(t, ErrMissingAuthorizationHeader.Code, w.Code)
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with invalid token format", func(t *testing.T) {
		t.Parallel()

		router := setupTestRouter(nil, zap.NewNop())
		req, w := newRequest(http.MethodGet, "/me", "invalid-format")
		router.ServeHTTP(w, req)

		require.Equal(t, ErrInvalidAuthorizationHeaderFormat.Code, w.Code)
		require.Equal(t, "GET, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		require.Contains(t, w.Body.String(), `"error":`)
	})

	t.Run("GET with empty token after Bearer", func(t *testing.T) {
		t.Parallel()

		router := setupTestRouter(nil, zap.NewNop())
		req, w := newRequest(http.MethodGet, "/me", "Bearer ")
		router.ServeHTTP(w, req)

		require.Equal(t, ErrEmptyBearerToken.Code, w.Code)
		require.Contains(t, w.Body.String(), `"error":`)
	})
}
