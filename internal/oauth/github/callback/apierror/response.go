package apierror

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error ErrorCode `json:"error"`
}

func Respond(c *gin.Context, apiErr *APIError) {
	_ = c.Error(apiErr)

	c.JSON(apiErr.HTTPStatus, gin.H{"error": apiErr.Code})
}
