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
