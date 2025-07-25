package info

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/config"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type InfoHandler struct {
	Logger *zap.Logger
}

func NewInfoHandler(logger *zap.Logger) *InfoHandler {
	return &InfoHandler{Logger: logger}
}

func (h *InfoHandler) Serve(c *gin.Context) {
	h.Logger.Info("/info access invoked")

	resp := InfoResponse{
		Message: "Welcome to IdP Proxy",
		OpenAPI: config.GetOpenAPIURL(),
	}
	c.JSON(http.StatusOK, resp)
}
