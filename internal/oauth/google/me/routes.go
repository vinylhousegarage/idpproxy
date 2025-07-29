package me

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r gin.IRoutes, googleDeps *deps.GoogleDependencies) {
	h := NewMeHandler(googleDeps.Verifier, googleDeps.Logger)
	r.GET("/me", h.Serve)
}
