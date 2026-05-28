package token

import "errors"

var (
	ErrEmptyBody          = errors.New("token: empty response body")
	ErrGitHubOAuthError   = errors.New("token: github oauth error")
	ErrMissingAccessToken = errors.New("token: missing access_token in response")
	ErrNon2xxStatus       = errors.New("token: non-2xx response")
	ErrParseFormBody      = errors.New("token: parse form body")
)
