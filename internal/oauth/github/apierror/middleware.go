package apierror

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
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
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			fields = append(fields,
				zap.String("code", string(apiErr.Code)),
				zap.Int("status", apiErr.HTTPStatus),
				zap.Error(apiErr.Err),
			)

			for i, info := range apiErr.Internal {
				key := fmt.Sprintf("detail_%d", i+1)
				fields = append(fields, zap.String(key, info))
			}

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
