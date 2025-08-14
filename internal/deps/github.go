package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

type GitHubOAuthDependencies struct {
	Config *config.GitHubOAuthConfig
	Logger *zap.Logger
}

func NewGitHubOAuthDeps(cfg *config.GitHubOAuthConfig, logger *zap.Logger) *GitHubOAuthDependencies {
	return &GitHubOAuthDependencies{
		Config: cfg,
		Logger: logger,
	}
}
