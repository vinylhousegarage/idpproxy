package me

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	handler := &GitHubMeHandler{}

	require.NotPanics(t, func() {
		RegisterRoutes(r, handler)
	})

	require.True(
		t,
		hasRegisteredRoute(r.Routes(), "GET", "/github/me"),
		"GET /github/me route was not registered",
	)
}

func hasRegisteredRoute(
	routes []gin.RouteInfo,
	method string,
	path string,
) bool {
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return true
		}
	}

	return false
}
