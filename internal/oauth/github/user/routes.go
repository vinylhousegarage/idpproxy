package user

import (
	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func RegisterRoutes(r gin.IRouter, githubAPIDeps *deps.GitHubAPIDependencies) {
	r.GET("/github/user", NewGitHubUserHandler(githubAPIDeps))
}
