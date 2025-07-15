package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthHandler struct {
	Logger *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{Logger: logger}
}

func (h *HealthHandler) Serve(c *gin.Context) {
	h.Logger.Info("health check invoked")

	resp := HealthResponse{Status: "healthy"}
	c.JSON(http.StatusOK, resp)
}
