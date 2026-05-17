package callback

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			logger.Error("request failed",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(err),
			)

			WriteError(c, err)
		}
	}
}
