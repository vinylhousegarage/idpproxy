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

	serverErr := ServerError(err)

	c.JSON(serverErr.HTTPStatus, ErrorResponse{
		Error: serverErr.Code,
	})
}
