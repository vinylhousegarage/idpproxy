package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func TestRouter_ErrorLoggerMiddlewareIsApplied(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.ErrorLevel)

	deps := RouterDeps{
		GitHubAPI:   &deps.GitHubAPIDependencies{},
		GitHubOAuth: &deps.GitHubOAuthDependencies{},
		Google:      &deps.GoogleDependencies{},
		Logger:      zap.New(core),
		System:      &deps.SystemDependencies{},
	}
	r := gin.New()
	RegisterRoutes(r, deps)

	r.GET("/test-error-middleware", func(c *gin.Context) {
		_ = c.Error(errors.New("test error"))
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-error-middleware", nil)
	r.ServeHTTP(w, req)

	if logs.Len() == 0 {
		t.Fatal("error log was not recorded")
	}
}
