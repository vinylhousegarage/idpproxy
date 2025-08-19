package login

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.Engine, githubOAuthDeps *deps.GitHubOAuthDependencies) {
	h := NewGitHubLoginHandler(githubOAuthDeps)
	r.GET("/github/login", h.Serve)
}
