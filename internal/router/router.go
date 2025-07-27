package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/loginfirebase"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/info"
)

func NewRouter(
	googleDeps *deps.GoogleDependencies,
	systemDeps *deps.SystemDependencies,
	publicFS http.FileSystem,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("login.html", publicFS)
	})

	r.StaticFS("/public", publicFS)

	googleGroup := r.Group("/google")
	loginfirebase.RegisterRoutes(googleGroup, googleDeps)

	health.RegisterRoutes(r, systemDeps)
	info.RegisterRoutes(r, systemDeps)

	return r
}
