package apierror

import "github.com/gin-gonic/gin"

func WriteError(c *gin.Context, apiErr *APIError) {
	defer c.Abort()

	c.JSON(apiErr.GetHTTPStatus(), ErrorResponse{
		Error: apiErr.Code,
	})
}
