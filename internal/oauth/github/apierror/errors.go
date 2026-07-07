package apierror

import "errors"

var (
	// callback
	ErrMissingGitHubCode  = errors.New(string(ErrorCodeMissingGitHubCode))
	ErrMissingState       = errors.New(string(ErrorCodeMissingState))
	ErrInvalidCookieState = errors.New(string(ErrorCodeInvalidCookieState))
	ErrInvalidQueryState  = errors.New(string(ErrorCodeInvalidQueryState))
	ErrInvalidState       = errors.New(string(ErrorCodeInvalidState))

	// token
	ErrBuildAccessTokenRequest  = errors.New(string(ErrorCodeBuildAccessTokenRequest))
	ErrGitHubAccessTokenRequest = errors.New(string(ErrorCodeGitHubAccessTokenRequest))
	ErrGitHubTokenRequest       = errors.New(string(ErrorCodeGitHubTokenRequest))
	ErrGitHubTokenExchange      = errors.New(string(ErrorCodeGitHubTokenExchange))

	// user
	ErrGitHubUserRequestBuild = errors.New(string(ErrorCodeGitHubUserRequestBuild))
	ErrGitHubUserRequest      = errors.New(string(ErrorCodeGitHubUserRequest))
	ErrGitHubUserDecode       = errors.New(string(ErrorCodeGitHubUserDecode))

	// internal
	ErrInternalServerError = errors.New(string(ErrorCodeInternalServerError))
	ErrProxyCodeIssue      = errors.New(string(ErrorCodeProxyCodeIssue))
	ErrUserUpsert          = errors.New(string(ErrorCodeUserUpsert))
)
