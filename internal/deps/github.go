package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
)

type GitHubDependencies struct {
	Config *config.GitHubConfig
	Logger *zap.Logger
}

func NewGitHubDeps(cfg *config.GitHubConfig, logger *zap.Logger) *GitHubDependencies {
	return &GitHubDependencies{
		Config: cfg,
		Logger: logger,
	}
}
