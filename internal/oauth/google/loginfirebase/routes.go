package loginfirebase

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r gin.IRoutes, googleDeps *deps.GoogleDependencies) {
	h := NewLoginFirebaseHandler(googleDeps.Verifier, googleDeps.Logger)
	r.POST("/google/login/firebase", h.Serve)
}
