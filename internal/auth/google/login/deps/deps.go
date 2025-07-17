package deps

import (
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
