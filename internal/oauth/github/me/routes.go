package me

import (
	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

func RegisterRoutes(
	r gin.IRouter,
	handler *GitHubMeHandler,
	validator session.SessionValidator,
) {
	r.GET(
		"/github/me",
		session.RequireSession(validator),
		handler.Serve,
	)
}
