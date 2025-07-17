package root

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, di *deps.Dependencies) {
	h := NewRootHandler(di)
	r.GET("/", h.Serve)
}
