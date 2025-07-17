package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/auth/google/login"
	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/root"
)

func NewRouter(
	logger *zap.Logger,
	metadataURL string,
	cfg config.GoogleConfig,
	client httpclient.HTTPClient,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	googleGroup := r.Group("google")
	login.RegisterRoutes(googleGroup, metadataURL, cfg, client, logger)

	systemGroup := r.Group("")
	health.RegisterRoutes(systemGroup, logger)
	root.RegisterRoutes(systemGroup, logger)

	return r
}
