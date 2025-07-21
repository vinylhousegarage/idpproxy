package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/login"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/root"
)

func NewRouter(di *deps.Dependencies, public http.FileSystem) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.StaticFS("/public", public)

	googleGroup := r.Group("google")
	login.RegisterRoutes(googleGroup, di)

	systemGroup := r.Group("")
	health.RegisterRoutes(systemGroup, di)
	root.RegisterRoutes(systemGroup, di)

	return r
}
