package router

import (
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/root"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	system := r.Group("")
	health.RegisterRoutes(system, logger)
	root.RegisterRoutes(system, logger)

	return r
}
