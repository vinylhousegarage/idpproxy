package main

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/repository"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/internal/server"
	"github.com/vinylhousegarage/idpproxy/public"
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

	ctx := context.Background()

	app, err := deps.NewFirebaseApp(ctx)
	if err != nil {
		logger.Fatal("failed to initialize Firebase App", zap.Error(err))
	}

	firestoreClient, err := deps.NewFirestoreClient(ctx, app, logger)
	if err != nil {
		logger.Fatal("failed to initialize Firestore client", zap.Error(err))
	}

	googleConfig, err := config.LoadGoogleConfig()
	if err != nil {
		logger.Fatal("failed to load google config", zap.Error(err))
	}

	metadataURL := config.GoogleOIDCMetadataURL
	httpClient := &http.Client{}
	tokenRepo := repository.NewGoogleTokenRepository(firestoreClient, logger)

	di := deps.New(metadataURL, googleConfig, httpClient, tokenRepo, logger)

	r := router.NewRouter(di, http.FS(public.PublicFS))

	server.StartServer(r, logger)
}
