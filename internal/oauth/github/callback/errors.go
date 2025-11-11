package callback

import "errors"

var (
	ErrEmptyBody          = errors.New("empty response body")
	ErrEncrypt            = errors.New("github: encrypt failed")
	ErrGitHubOAuthError   = errors.New("github oauth error")
	ErrInvalidInput       = errors.New("github: invalid input")
	ErrMissingAccessToken = errors.New("missing access_token in response")
	ErrNon2xxStatus       = errors.New("non-2xx response")
	ErrParseFormBody      = errors.New("parse form body")
	ErrPersist            = errors.New("github: persist failed")
)
