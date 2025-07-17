package login

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	metadataURL string,
	cfg config.GoogleConfig,
	client httpclient.HTTPClient,
	logger *zap.Logger,
) {
	h := NewGoogleLoginHandler(metadataURL, cfg, client, logger)
	r.GET("/google/login", h.Serve)
}
