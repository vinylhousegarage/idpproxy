package login

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func mockGoogleConfig() config.GoogleConfig {
	return config.GoogleConfig{
		ClientID:     "test-client-id",
		ClientSecret: "secret",
		RedirectURI:  "https://idpproxy.com/callback",
		ResponseType: "code",
		Scope:        "openid email",
		AccessType:   "offline",
		Prompt:       "consent",
	}
}

var mockMeta = `{"authorization_endpoint": "https://accounts.google.com/o/oauth2/auth"}`

func TestGoogleLoginHandler_Serve_Success(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := mockGoogleConfig()
	logger := zaptest.NewLogger(t)

	client := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(strings.NewReader(mockMeta))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}, nil
		},
	}

	handler := NewGoogleLoginHandler(mockMeta, cfg, client, logger)
	router.GET("/google/login", handler.Serve)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/google/login", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	location := w.Header().Get("Location")
	assert.Contains(t, location, "https://accounts.google.com/o/oauth2/auth")
	assert.Contains(t, location, "client_id=test-client-id")
	assert.Contains(t, w.Header().Get("Set-Cookie"), "oauth_state=")
}

func TestGoogleLoginHandler_Serve_MetadataError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := mockGoogleConfig()
	logger := zaptest.NewLogger(t)

	client := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("failed to fetch metadata")
		},
	}

	handler := NewGoogleLoginHandler(mockMeta, cfg, client, logger)
	router.GET("/google/login", handler.Serve)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/google/login", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
