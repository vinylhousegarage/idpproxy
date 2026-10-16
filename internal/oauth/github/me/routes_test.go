package me

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

type fakeSessionValidator struct{}

func (fakeSessionValidator) Validate(
	_ context.Context,
	_ string,
) (*authsession.Session, error) {
	return nil, nil
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	handler := &GitHubMeHandler{}
	validator := fakeSessionValidator{}

	require.NotPanics(t, func() {
		RegisterRoutes(r, handler, validator)
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
