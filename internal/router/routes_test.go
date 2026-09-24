package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
)

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	t.Run("registers_GitHub_callback_route", func(t *testing.T) {
		t.Parallel()

		gin.SetMode(gin.TestMode)

		r := gin.New()
		d := RouterDeps{
			GitHubAPI:      &deps.GitHubAPIDependencies{},
			GitHubOAuth:    &deps.GitHubOAuthDependencies{},
			GitHubCallback: &callback.GitHubCallbackHandler{},
			Google:         &deps.GoogleDependencies{},
			Logger:         zap.NewNop(),
			System:         &deps.SystemDependencies{},
		}

		require.NotPanics(t, func() {
			RegisterRoutes(r, d)
		})

		require.True(
			t,
			hasRoute(r.Routes(), "GET", "/oauth/github/callback"),
			"GitHub callback route was not registered",
		)
	})

	t.Run("panics_when_GitHub_callback_handler_is_missing", func(t *testing.T) {
		t.Parallel()

		gin.SetMode(gin.TestMode)

		r := gin.New()
		d := RouterDeps{
			GitHubAPI:   &deps.GitHubAPIDependencies{},
			GitHubOAuth: &deps.GitHubOAuthDependencies{},
			Google:      &deps.GoogleDependencies{},
			Logger:      zap.NewNop(),
			System:      &deps.SystemDependencies{},
		}

		require.PanicsWithValue(
			t,
			"router: missing dependencies",
			func() {
				RegisterRoutes(r, d)
			},
		)
	})
}

func hasRoute(
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
