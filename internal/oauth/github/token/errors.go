package token

import "errors"

var (
	ErrEmptyBody          = errors.New("callback: empty response body")
	ErrGitHubOAuthError   = errors.New("callback: github oauth error")
	ErrMissingAccessToken = errors.New("callback: missing access_token in response")
	ErrNon2xxStatus       = errors.New("callback: non-2xx response")
	ErrParseFormBody      = errors.New("callback: parse form body")
)
