package callback

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
	"go.uber.org/zap"
)

func ErrorLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		fields := []zap.Field{
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		}

		logLevelFunc := logger.Error
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			fields = append(fields,
				zap.String("code", string(apiErr.Code)),
				zap.Int("status", apiErr.HTTPStatus),
				zap.String("internal_info", apiErr.Internal),
				zap.Error(apiErr.Err),
			)
			if apiErr.HTTPStatus >= 400 && apiErr.HTTPStatus < 500 {
				logLevelFunc = logger.Warn
			}
		} else {
			fields = append(fields, zap.Error(err))
		}

		logLevelFunc("request failed", fields...)

		WriteError(c, err)
	}
}
