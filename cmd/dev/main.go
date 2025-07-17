package main

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
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

	cfg, err := config.LoadGoogleConfig()
	if err != nil {
		panic("failed to initialize cfg: " + err.Error())
	}

	di := deps.New(cfg, logger)

	r := router.NewRouter(di)

	server.StartServer(r, logger)
}
