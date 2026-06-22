package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRouter_ErrorLoggerMiddlewareIsApplied(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zap.ErrorLevel)
	testLogger := zap.New(core)

	deps := router.RouterDeps{
		Logger: testLogger,
	}
	r := gin.New()
	router.RegisterRoutes(r, deps)

	r.GET("/test-error-middleware", func(c *gin.Context) {
		c.Error(errors.New("test error"))
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-error-middleware", nil)
	r.ServeHTTP(w, req)

	if logs.Len() == 0 {
		t.Fatal("error log was not recorded")
	}
}
