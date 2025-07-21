package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/repository"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/internal/server"
)


func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	googleConfig, err := config.LoadGoogleConfig()
	if err != nil {
		logger.Fatal("failed to load google config", zap.Error(err))
	}

	firestoreConfig := config.LoadFirestoreConfig()
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, firestoreConfig.ProjectID)
	if err != nil {
		logger.Fatal("failed to initialize Firestore client", zap.Error(err))
	}

	metadataURL := config.GoogleOIDCMetadataURL
	httpClient := &http.Client{}
	tokenRepo := repository.NewGoogleTokenRepository(firestoreClient, logger)

	di := deps.New(metadataURL, googleConfig, httpClient, tokenRepo, logger)

	r := router.NewRouter(di)

	server.StartServer(r, logger)
}
