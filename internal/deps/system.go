package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

type SystemDependencies struct {
	MetadataURL string
	HTTPClient  httpclient.HTTPClient
	Logger      *zap.Logger
}

func NewSystemDeps(
	metadataURL string,
	httpClient httpclient.HTTPClient,
	logger *zap.Logger,
) *SystemDependencies {
	return &SystemDependencies{
		MetadataURL: metadataURL,
		HTTPClient:  httpClient,
		Logger:      logger,
	}
}
