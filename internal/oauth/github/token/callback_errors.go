package token

import "github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"

const (
	ErrorCodeBuildAccessTokenRequest apierror.ErrorCode = "build_access_token_request_failed"
	ErrorCodeGitHubTokenRequest      apierror.ErrorCode = "github_token_request_failed"
	ErrorCodeGitHubTokenExchange     apierror.ErrorCode = "github_token_exchange_failed"
)
