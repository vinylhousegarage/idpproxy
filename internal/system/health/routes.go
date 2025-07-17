package health

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, di *deps.Dependencies) {
	h := NewHealthHandler(di)
	r.GET("/health", h.Serve)
}
