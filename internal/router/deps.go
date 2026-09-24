package router

import (
	"io/fs"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
)

type RouterDeps struct {
	FS             fs.FS
	GitHubAPI      *deps.GitHubAPIDependencies
	GitHubOAuth    *deps.GitHubOAuthDependencies
	GitHubCallback *callback.GitHubCallbackHandler
	Google         *deps.GoogleDependencies
	Logger         *zap.Logger
	System         *deps.SystemDependencies
}

func NewRouterDeps(
	fsys fs.FS,
	githubAPI *deps.GitHubAPIDependencies,
	githubOAuth *deps.GitHubOAuthDependencies,
	githubCallback *callback.GitHubCallbackHandler,
	google *deps.GoogleDependencies,
	logger *zap.Logger,
	system *deps.SystemDependencies,
) RouterDeps {
	return RouterDeps{
		FS:             fsys,
		GitHubAPI:      githubAPI,
		GitHubOAuth:    githubOAuth,
		GitHubCallback: githubCallback,
		Google:         google,
		Logger:         logger,
		System:         system,
	}
}
