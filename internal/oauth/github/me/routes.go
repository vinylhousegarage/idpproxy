package me

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r gin.IRouter,
	handler *GitHubMeHandler,
) {
	r.GET("/github/me", handler.Serve)
}
