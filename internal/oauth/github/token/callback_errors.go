package token

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

const (
	ErrorCodeBuildAccessTokenRequest  apierror.ErrorCode = "build_access_token_request_failed"
	ErrorCodeGitHubAccessTokenRequest apierror.ErrorCode = "github_access_token_request_failed"
	ErrorCodeGitHubTokenRequest       apierror.ErrorCode = "github_token_request_failed"
	ErrorCodeGitHubTokenExchange      apierror.ErrorCode = "github_token_exchange_failed"
)

func GitHubAccessTokenRequestError(err error, internal ...string) *apierror.APIError {
	return apierror.New(ErrorCodeGitHubAccessTokenRequest, http.StatusBadGateway, err, internal...)
}
