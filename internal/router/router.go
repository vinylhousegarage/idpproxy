package router

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(d RouterDeps) *gin.Engine {
	r := gin.New()

	if err := r.SetTrustedProxies(nil); err != nil {
		panic("router: failed to set trusted proxies: " + err.Error())
	}

	r.Use(gin.Recovery())

	RegisterRoutes(r, d)

	return r
}
