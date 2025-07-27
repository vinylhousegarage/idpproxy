package loginfirebase

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.RouterGroup, googleDeps *deps.GoogleDependencies) {
	h := NewLoginFirebaseHandler(googleDeps.Logger)
	r.GET("/loginfirebase", h.Serve)
}
