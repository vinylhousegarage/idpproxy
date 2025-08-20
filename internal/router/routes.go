package router

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/login"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/loginfirebase"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/me"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/info"
)

func RegisterRoutes(r *gin.Engine, d RouterDeps) {
	if d.GitHubAPI == nil || d.GitHubOAuth == nil || d.Google == nil || d.System == nil {
			panic("router: missing dependencies")
	}

	if d.FS != nil {
		r.GET("/",        func(c *gin.Context) { c.FileFromFS("login.html",   http.FS(d.FS)) })
		r.GET("/privacy", func(c *gin.Context) { c.FileFromFS("privacy.html", http.FS(d.FS)) })
		r.GET("/terms",   func(c *gin.Context) { c.FileFromFS("terms.html",   http.FS(d.FS)) })
	}

	// GitHub
	login.RegisterRoutes(r, d.GitHubOAuth)
	user.RegisterRoutes(r, d.GitHubAPI)

	// Google
	loginfirebase.RegisterRoutes(r, d.Google)
	me.RegisterRoutes(r, d.Google)

	// System
	health.RegisterRoutes(r, d.System)
	info.RegisterRoutes(r, d.System)
}
