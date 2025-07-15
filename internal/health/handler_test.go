package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestHealthHandler_Returns200(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := NewHealthHandler(zap.NewNop())
	router.GET("/health", handler.Serve)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"status":"healthy"}`, w.Body.String())
}
