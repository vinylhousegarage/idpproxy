package health

import (
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func RegisterRoutes(r *gin.RouterGroup, logger *zap.Logger) {
	h := NewHealthHandler(logger)
	r.GET("/health", h.Serve)
}
