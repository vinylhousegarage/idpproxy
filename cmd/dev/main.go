package main

import (
	"context"
	"net/http"
	"time"

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
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()
	firebaseCfg, err := config.LoadFirebaseConfig()
	if err != nil {
		logger.Fatal("failed to load Firebase config", zap.Error(err))
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(firebaseCfg.CredentialsJSON))
	if err != nil {
		logger.Fatal("failed to initialize Firebase App", zap.Error(err))
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Fatal("failed to initialize Firebase Auth client", zap.Error(err))
	}

	googleDeps := deps.NewGoogleDeps(authClient, logger)

	githubCfg, err := config.LoadGitHubDevOAuthConfig()
	if err != nil {
		logger.Fatal("failed to load GitHub config", zap.Error(err))
	}

	githubOAuthDeps := deps.NewGitHubOAuthDeps(githubCfg, logger)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	githubAPICfg := config.LoadGitHubAPIConfig()
	githubAPIDeps := deps.NewGitHubAPIDeps(githubAPICfg, httpClient, logger)

	systemDeps := deps.NewSystemDeps(config.GoogleOIDCMetadataURL, httpClient, logger)

	d := router.NewRouterDeps(public.PublicFS, githubAPIDeps, githubOAuthDeps, googleDeps, systemDeps)
	r := router.NewRouter(d)

	logger.Info("starting idpproxy (dev)", zap.String("addr", ":"+config.GetPort()))
	server.StartServer(r, logger)
}
