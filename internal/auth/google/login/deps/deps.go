package deps

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

type Dependencies struct {
	MetadataURL string
	Config      config.GoogleConfig
	HTTPClient  httpclient.HTTPClient
	Logger      *zap.Logger
}

func New(cfg config.GoogleConfig, logger *zap.Logger) *Dependencies {
	return &Dependencies{
		MetadataURL: config.GoogleOIDCMetadataURL,
		Config:      cfg,
		HTTPClient:  &http.Client{},
		Logger:      logger,
	}
}
