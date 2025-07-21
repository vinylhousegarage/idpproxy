package router

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/login"
	"github.com/vinylhousegarage/idpproxy/internal/system/health"
	"github.com/vinylhousegarage/idpproxy/internal/system/root"
)

//go:embed ../../public/*
var publicFS embed.FS

func NewRouter(di *deps.Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.StaticFS("/public", http.FS(publicFS))

	googleGroup := r.Group("google")
	login.RegisterRoutes(googleGroup, di)

	systemGroup := r.Group("")
	health.RegisterRoutes(systemGroup, di)
	root.RegisterRoutes(systemGroup, di)

	return r
}
