package router

import (
	"io/fs"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	githubOAuthDeps *deps.GitHubOAuthDependencies,
	githubAPIDeps *deps.GitHubAPIDependencies,
	googleDeps *deps.GoogleDependencies,
	systemDeps *deps.SystemDependencies,
	publicFS http.FileSystem,
) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())

	RegisterRoutes(r, d)

	return r
}
