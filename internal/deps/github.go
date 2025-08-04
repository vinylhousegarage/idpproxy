package deps

import (
	"go.uber.org/zap"
)

type GitHubDependencies struct {
	Logger *zap.Logger
}

func NewGitHubDeps(logger *zap.Logger) *GitHubDependencies {
	return &GitHubDependencies{
		Logger: logger,
	}
}
