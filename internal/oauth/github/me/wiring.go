package me

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
)

type GitHubMeService interface {
	Get(
		ctx context.Context,
		firebaseUID string,
	) (*response.GitHubUserAPIResponse, error)
}

type GitHubMeHandler struct {
	Service GitHubMeService
}

func NewGitHubMeHandler(
	service GitHubMeService,
) *GitHubMeHandler {
	return &GitHubMeHandler{
		Service: service,
	}
}
