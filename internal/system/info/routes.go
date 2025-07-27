package info

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, systemDeps *deps.SystemDependencies) {
	h := NewInfoHandler(systemDeps.Logger)
	r.GET("/info", h.Serve)
}
