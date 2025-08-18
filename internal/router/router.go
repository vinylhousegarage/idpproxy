package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/login"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/loginfirebase"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/me"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/info"
)

func NewRouter(
	githubDeps *deps.GitHubOAuthDependencies,
	githubAPIDeps *deps.GitHubAPIDependencies,
	googleDeps *deps.GoogleDependencies,
	systemDeps *deps.SystemDependencies,
	publicFS http.FileSystem,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("login.html", publicFS)
	})

	r.GET("/privacy", func(c *gin.Context) {
		c.FileFromFS("privacy.html", publicFS)
	})

	r.GET("/terms", func(c *gin.Context) {
		c.FileFromFS("terms.html", publicFS)
	})

	r.StaticFS("/public", publicFS)

	login.RegisterRoutes(r, githubDeps)
	user.RegisterRoutes(r, githubAPIDeps)

	loginfirebase.RegisterRoutes(r, googleDeps)
	me.RegisterRoutes(r, googleDeps)

	health.RegisterRoutes(r, systemDeps)
	info.RegisterRoutes(r, systemDeps)

	return r
}
