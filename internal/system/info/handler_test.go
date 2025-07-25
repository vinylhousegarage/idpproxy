package info

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestNewInfoHandler_Returns200JSON(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewInfoHandler(zap.NewNop())
	router.GET("/info", handler.Serve)

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{
		"message": "Welcome to IdP Proxy",
		"openapi": "http://localhost:9000/openapi.json"
	}`, w.Body.String())
}
