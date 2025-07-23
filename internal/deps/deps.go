package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

type Dependencies struct {
	MetadataURL string
	HTTPClient  httpclient.HTTPClient
	Logger      *zap.Logger
}

func New(
	metadataURL string,
	httpClient httpclient.HTTPClient,
	logger *zap.Logger,
) *Dependencies {
	return &Dependencies{
		MetadataURL: metadataURL,
		HTTPClient:  httpClient,
		Logger:      logger,
	}
}
