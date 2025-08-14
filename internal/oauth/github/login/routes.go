package login

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r *gin.Engine, githubDeps *deps.GitHubOAuthDependencies) {
	h := NewGitHubLoginHandler(githubDeps)
	r.GET("/github/login", h.Serve)
}
