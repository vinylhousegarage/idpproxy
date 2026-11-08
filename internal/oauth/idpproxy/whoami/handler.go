package whoami

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

type WhoAmIHandler struct{}

func (h *WhoAmIHandler) Serve(c *gin.Context) {
	userID, ok := c.Get(authsession.ContextUserIDKey)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}
