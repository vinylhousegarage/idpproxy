package deps

import (
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
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

type GitHubAPIDependencies struct {
	APIVersion string
	BaseURL    string
	HTTPClient httpclient.HTTPClient
	Logger     *zap.Logger
	UserAgent  string
}

func NewGitHubAPIDeps(
	cfg *config.GitHubAPIConfig,
	client httpclient.HTTPClient,
	logger *zap.Logger,
) *GitHubAPIDependencies {
	return &GitHubAPIDependencies{
		APIVersion: cfg.APIVersion,
		BaseURL:    cfg.BaseURL,
		HTTPClient: client,
		Logger:     logger,
		UserAgent:  cfg.UserAgent,
	}
}
