package health

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, di *deps.Dependencies) {
	h := NewHealthHandler(di.Logger)
	r.GET("/health", h.Serve)
}
