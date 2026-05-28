package token

import (
	"errors"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

const (
	ErrorCodeBuildAccessTokenRequest apierror.ErrorCode = "build_access_token_request_failed"
	ErrorCodeGitHubTokenRequest      apierror.ErrorCode = "github_token_request_failed"
	ErrorCodeGitHubTokenExchange     apierror.ErrorCode = "github_token_exchange_failed"
)

var (
	ErrEmptyBody          = errors.New("token: empty response body")
	ErrGitHubOAuthError   = errors.New("token: github oauth error")
	ErrMissingAccessToken = errors.New("token: missing access_token in response")
	ErrNon2xxStatus       = errors.New("token: non-2xx response")
	ErrParseFormBody      = errors.New("token: parse form body")
)
