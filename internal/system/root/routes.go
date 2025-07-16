package root

import (
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func RegisterRoutes(r *gin.RouterGroup, logger *zap.Logger) {
	h := NewRootHandler(logger)
	r.GET("/", h.Serve)
}
