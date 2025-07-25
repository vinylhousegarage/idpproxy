package info

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, di *deps.Dependencies) {
	h := NewInfoHandler(di.Logger)
	r.GET("/info", h.Serve)
}
