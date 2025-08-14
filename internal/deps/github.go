package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

type GitHubOAuthDependencies struct {
	Config *config.GitHubConfig
	Logger *zap.Logger
}

func NewGitHubOAuthDeps(cfg *config.GitHubConfig, logger *zap.Logger) *GitHubOAuthDependencies {
	return &GitHubOAuthDependencies{
		Config: cfg,
		Logger: logger,
	}
}
