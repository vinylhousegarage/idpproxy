package routertest

import (
	"go.uber.org/zap"

	githubcallback "github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
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

	return router.NewRouterDeps(
		public.PublicFS,
		testhelpers.NewMockGitHubAPIDeps(logger),
		testhelpers.NewMockGitHubOAuthDeps(logger),
		&githubcallback.GitHubCallbackHandler{},
		&githubme.GitHubMeHandler{},
		testhelpers.NewMockSessionValidator(),
		&idpproxytoken.Handler{},
		&whoami.WhoAmIHandler{},
		testhelpers.NewMockGoogleDeps(logger),
		logger,
		testhelpers.NewMockSystemDeps(logger),
	)
}
