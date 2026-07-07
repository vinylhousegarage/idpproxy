package apierror

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	defer c.Abort()

	var apiErr *APIError

	if errors.As(err, &apiErr) {
		c.JSON(apiErr.HTTPStatus, ErrorResponse{
			Error: apiErr.Code,
		})

		return
	}

	internalServerErr := InternalServerError(err)

	c.JSON(internalServerErr.HTTPStatus, ErrorResponse{
		Error: internalServerErr.Code,
	})
}
