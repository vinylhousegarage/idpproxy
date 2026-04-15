package callback

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"
)

func WriteError(c *gin.Context, err error) {
	var apiErr *apierror.APIError

	if errors.As(err, &apiErr) {
		c.JSON(apiErr.HTTPStatus, apierror.ErrorResponse{
			Error: apiErr.Code,
		})

		return
	}

	internalErr := apierror.Internal(err)

	c.JSON(internalErr.HTTPStatus, apierror.ErrorResponse{
		Error: internalErr.Code,
	})
}
