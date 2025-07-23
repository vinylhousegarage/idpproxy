package main

import (
	"net/http"

	"go.uber.org/zap"

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

	di := deps.New(metadataURL, httpClient, logger)

	r := router.NewRouter(di, http.FS(public.PublicFS))

	server.StartServer(r, logger)
}