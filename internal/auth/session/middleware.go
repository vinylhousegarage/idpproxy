package session

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName = "idpproxy_session"
	ContextUserIDKey  = "current_user_id"
)

type SessionValidator interface {
	Validate(
		ctx context.Context,
		sessionID string,
	) (*Session, error)
}

func RequireSession(validator SessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validator == nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "server not ready"},
			)

			return
		}

		cookie, err := c.Request.Cookie(SessionCookieName)
		if err != nil || cookie == nil || cookie.Value == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "unauthorized"},
			)

			return
		}

		s, err := validator.Validate(c.Request.Context(), cookie.Value)
		if err != nil || s == nil || s.UserID == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "unauthorized"},
			)

			return
		}

		c.Set(ContextUserIDKey, s.UserID)
		c.Next()
	}
}
