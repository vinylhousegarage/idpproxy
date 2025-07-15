package main

import (
	"idpproxy/internal/router"
	"idpproxy/internal/server"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	r := router.NewRouter(logger)

	server.StartServer(r, logger)
}
