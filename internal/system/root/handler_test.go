package root

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestNewRootHandler_Returns200JSON(t *testing.T) {
	t.Setenv("OPENAPI_URL", "https://test.example.com/openapi.json")

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewRootHandler(zap.NewNop())
	router.GET("/", handler.Serve)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{
		"message": "Welcome to IdP Proxy",
		"openapi": "https://test.example.com/openapi.json"
	}`, w.Body.String())
}
