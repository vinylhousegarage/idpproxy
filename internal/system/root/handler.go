package root

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/config"

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
		OpenAPI: config.GetOpenAPIURL(),
	}
	c.JSON(http.StatusOK, resp)
}
