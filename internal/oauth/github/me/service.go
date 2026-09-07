package me

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
	githubstore "github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
	githubuser "github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
)

type GitHubTokenRepository interface {
	GetByFirebaseUID(
		ctx context.Context,
		firebaseUID string,
	) (*githubstore.GitHubTokenRecord, error)

	TouchLastUsed(
		ctx context.Context,
		firebaseUID string,
	) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	tokenRepo  GitHubTokenRepository
	httpClient HTTPClient
}

func NewService(
	tokenRepo GitHubTokenRepository,
	httpClient HTTPClient,
) *Service {
	return &Service{
		tokenRepo:  tokenRepo,
		httpClient: httpClient,
	}
}

func (s *Service) Get(
	ctx context.Context,
	firebaseUID string,
) (*response.GitHubUserAPIResponse, error) {
	token, err := s.tokenRepo.GetByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("get GitHub token: %w", err)
	}

	req, err := githubuser.NewGitHubUserRequest(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("build GitHub user request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request GitHub user: %w", err)
	}

	githubUser, err := githubuser.DecodeGitHubUserResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("decode GitHub user response: %w", err)
	}

	if err := s.tokenRepo.TouchLastUsed(ctx, firebaseUID); err != nil {
		return nil, fmt.Errorf("touch GitHub token last used: %w", err)
	}

	return githubUser, nil
}
