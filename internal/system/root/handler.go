package root

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type RootHandler struct {
	Logger *zap.Logger
}

func NewRootHandler(logger *zap.Logger) *RootHandler {
	return &RootHandler{Logger: logger}
}

func (h *RootHandler) Serve(c *gin.Context) {
	h.Logger.Info("/ access invoked")

	resp := RootResponse{
		Message: "Welcome to IdP Proxy",
		OpenAPI: os.Getenv("OPENAPI_URL"),
	}
	c.JSON(http.StatusOK, resp)
}
