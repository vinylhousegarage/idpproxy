package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/repository"
)

type Dependencies struct {
	MetadataURL     string
	Config          *config.GoogleConfig
	HTTPClient      httpclient.HTTPClient
	FirestoreClient repository.GoogleTokenStore
	Logger          *zap.Logger
}

func New(
	metadataURL string,
	googleConfig *config.GoogleConfig,
	httpClient httpclient.HTTPClient,
	firestoreClient repository.GoogleTokenStore,
	logger *zap.Logger,
) *Dependencies {
	return &Dependencies{
		MetadataURL:     metadataURL,
		Config:          googleConfig,
		HTTPClient:      httpClient,
		FirestoreClient: firestoreClient,
		Logger:          logger,
	}
}
