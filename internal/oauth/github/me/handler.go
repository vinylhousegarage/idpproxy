package me

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

func (h *GitHubMeHandler) Serve(c *gin.Context) {
	if h == nil || h.Service == nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "server not ready"},
		)

		return
	}

	value, ok := c.Get(session.ContextUserIDKey)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)

		return
	}

	firebaseUID, ok := value.(string)
	if !ok || firebaseUID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)

		return
	}

	githubUser, err := h.Service.Get(c.Request.Context(), firebaseUID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get GitHub user"},
		)

		return
	}

	c.JSON(http.StatusOK, githubUser)
}
