package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/root"
)

func NewRouter(di *deps.Dependencies, publicFS http.FileSystem) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.StaticFS("/public", publicFS)

	systemGroup := r.Group("")
	health.RegisterRoutes(systemGroup, di)
	root.RegisterRoutes(systemGroup, di)

	return r
}
