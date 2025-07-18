package testhelpers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func NewMockDeps(logger *zap.Logger) *deps.Dependencies {
	return &deps.Dependencies{
		MetadataURL: "https://accounts.google.com/.well-known/openid-configuration",
		Config: &config.GoogleConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURI:  "https://localhost/callback",
			ResponseType: "code",
			Scope:        "openid email profile",
			AccessType:   "offline",
			Prompt:       "consent",
		},
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}
}
