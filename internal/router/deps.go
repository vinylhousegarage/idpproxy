package router

import (
	"io/fs"

	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	idpproxytoken "github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/token"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/whoami"
)

type RouterDeps struct {
	FS               fs.FS
	GitHubAPI        *deps.GitHubAPIDependencies
	GitHubOAuth      *deps.GitHubOAuthDependencies
	GitHubCallback   *callback.GitHubCallbackHandler
	GitHubMe         *githubme.GitHubMeHandler
	SessionValidator authsession.SessionValidator
	IDPProxyToken    *idpproxytoken.Handler
	WhoAmI           *whoami.WhoAmIHandler
	Google           *deps.GoogleDependencies
	Logger           *zap.Logger
	System           *deps.SystemDependencies
}
