package loginfirebase

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/verify"
)

func RegisterRoutes(r *gin.RouterGroup, verifier verify.Verifier, logger *zap.Logger) {
	h := NewLoginFirebaseHandler(verifier, logger)
	r.GET("/loginfirebase", h.Serve)
}
