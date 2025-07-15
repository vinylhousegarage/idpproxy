package server

import (
	"idpproxy/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StartServer(r *gin.Engine, logger *zap.Logger) {
	port := config.GetPort()

	logger.Info("Starting server", zap.String("port", port))

	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
