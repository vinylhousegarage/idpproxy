package router

import (
	"io/fs"

	"github.com/gin-gonic/gin"
)

func NewRouter(d RouterDeps) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())

	RegisterRoutes(r, d)

	return r
}
