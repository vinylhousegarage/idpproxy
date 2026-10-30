package router

import (
	"io/fs"

	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	idpproxytoken "github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/token"
)

type RouterDeps struct {
	FS               fs.FS
	GitHubAPI        *deps.GitHubAPIDependencies
	GitHubOAuth      *deps.GitHubOAuthDependencies
	GitHubCallback   *callback.GitHubCallbackHandler
	GitHubMe         *githubme.GitHubMeHandler
	SessionValidator authsession.SessionValidator
	IDPProxyToken    *idpproxytoken.Handler
	Google           *deps.GoogleDependencies
	Logger           *zap.Logger
	System           *deps.SystemDependencies
}

func NewRouterDeps(
	fsys fs.FS,
	githubAPI *deps.GitHubAPIDependencies,
	githubOAuth *deps.GitHubOAuthDependencies,
	githubCallback *callback.GitHubCallbackHandler,
	githubMe *githubme.GitHubMeHandler,
	sessionValidator authsession.SessionValidator,
	idpProxyToken *idpproxytoken.Handler,
	google *deps.GoogleDependencies,
	logger *zap.Logger,
	system *deps.SystemDependencies,
) RouterDeps {
	return RouterDeps{
		FS:               fsys,
		GitHubAPI:        githubAPI,
		GitHubOAuth:      githubOAuth,
		GitHubCallback:   githubCallback,
		GitHubMe:         githubMe,
		SessionValidator: sessionValidator,
		IDPProxyToken:    idpProxyToken,
		Google:           google,
		Logger:           logger,
		System:           system,
	}
}
