package callback

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter, handler *GitHubCallbackHandler) {
	r.GET("/oauth/github/callback", handler.Serve)
}
