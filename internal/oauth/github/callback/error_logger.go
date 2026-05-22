package callback

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"
)

func ErrorLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			fields := []zap.Field{
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			}

			if apiErr, ok := err.(*apierror.APIError); ok {
				fields = append(fields,
					zap.String("code", string(apiErr.Code)),
					zap.Int("status", apiErr.HTTPStatus),
					zap.String("internal_info", apiErr.Internal),
					zap.Error(apiErr.Err),
				)
			} else {
				fields = append(fields, zap.Error(err))
			}

			logger.Error("request failed", fields...)

			WriteError(c, err)
		}
	}
}
