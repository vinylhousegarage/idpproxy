package health

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.IRoutes, systemDeps *deps.SystemDependencies) {
	h := NewHealthHandler(systemDeps.Logger)
	r.GET("/health", h.Serve)
}
