package testhelpers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func NewMockDeps(logger *zap.Logger) *deps.Dependencies {
	return &deps.Dependencies{
		MetadataURL: "https://accounts.google.com/.well-known/openid-configuration",
		HTTPClient:  http.DefaultClient,
		Logger:      logger,
	}
}
