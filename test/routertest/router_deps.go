package routertest

import (
	"go.uber.org/zap"

	githubcallback "github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/logout"
	idpproxytoken "github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/token"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/whoami"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func NewMockRouterDeps(logger *zap.Logger) router.RouterDeps {
	if logger == nil {
		logger = zap.NewNop()
	}

	return router.RouterDeps{
		FS:               public.PublicFS,
		GitHubAPI:        testhelpers.NewMockGitHubAPIDeps(logger),
		GitHubOAuth:      testhelpers.NewMockGitHubOAuthDeps(logger),
		GitHubCallback:   &githubcallback.GitHubCallbackHandler{},
		GitHubMe:         &githubme.GitHubMeHandler{},
		SessionValidator: testhelpers.NewMockSessionValidator(),
		IDPProxyToken:    &idpproxytoken.Handler{},
		WhoAmI:           &whoami.WhoAmIHandler{},
		Logout:           logout.NewHandler(testhelpers.NewMockSessionInvalidator()),
		Google:           testhelpers.NewMockGoogleDeps(logger),
		Logger:           logger,
		System:           testhelpers.NewMockSystemDeps(logger),
	}
}
