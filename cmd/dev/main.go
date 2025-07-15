package main

import (
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/internal/server"

	"go.uber.org/zap"
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

	r := router.NewRouter(logger)

	server.StartServer(r, logger)
}
