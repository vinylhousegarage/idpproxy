package router

import (
	"io/fs"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

type RouterDeps struct {
	FS          fs.FS
	GitHubAPI   *deps.GitHubAPIDependencies
	GitHubOAuth *deps.GitHubOAuthDependencies
	Google      *deps.GoogleDependencies
	System      *deps.SystemDependencies
}

func NewRouterDeps(
	fsys        fs.FS,
	githubAPI   *deps.GitHubAPIDependencies,
	githubOAuth *deps.GitHubOAuthDependencies,
	google      *deps.GoogleDependencies,
	system      *deps.SystemDependencies,
) RouterDeps {
	return RouterDeps{
		FS:          fsys,
		GitHubAPI:   githubAPI,
		GitHubOAuth: githubOAuth,
		Google:      google,
		System:      system,
	}
}
