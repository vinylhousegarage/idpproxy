package main

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
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

	metadataURL := config.GoogleOIDCMetadataURL
	httpClient := &http.Client{}
	systemDeps := deps.NewSystemDeps(metadataURL, httpClient, logger)

	ctx := context.Background()
	firebaseCfg, err := config.LoadFirebaseConfig()
	if err != nil {
		logger.Fatal("failed to load Firebase config", zap.Error(err))
	}

	opt := option.WithCredentialsJSON(firebaseCfg.CredentialsJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Fatal("failed to initialize Firebase App", zap.Error(err))
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Fatal("failed to initialize Firebase Auth client", zap.Error(err))
	}

	googleDeps := deps.NewGoogleDeps(authClient, logger)

	githubCfg, err := config.LoadGitHubDevConfig()
	if err != nil {
		logger.Fatal("failed to load GitHub config", zap.Error(err))
	}

	githubDeps := deps.NewGitHubOAuthDeps(githubCfg, logger)

	r := router.NewRouter(githubDeps, googleDeps, systemDeps, http.FS(public.PublicFS))
	server.StartServer(r, logger)
}
