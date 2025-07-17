package login

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, di *deps.Dependencies) {
	h := NewGoogleLoginHandler(di)
	r.GET("/login", h.Serve)
}
