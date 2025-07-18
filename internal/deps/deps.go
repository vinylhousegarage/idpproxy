package deps

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

type Dependencies struct {
	MetadataURL     string
	Config          *config.GoogleConfig
	HTTPClient      httpclient.HTTPClient
	FirestoreClient *firestore.Client
	Logger          *zap.Logger
}

func New(
	metadataURL string,
	googleConfig *config.GoogleConfig,
	httpClient httpclient.HTTPClient,
	firestoreClient *firestore.Client,
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
