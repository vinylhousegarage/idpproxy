package callback

import "errors"

// HTTP/Parse layer (GitHub OAuth response handling)
var (
	ErrEmptyBody          = errors.New("callback: empty response body")
	ErrGitHubOAuthError   = errors.New("callback: github oauth error")
	ErrMissingAccessToken = errors.New("callback: missing access_token in response")
	ErrNon2xxStatus       = errors.New("callback: non-2xx response")
	ErrParseFormBody      = errors.New("callback: parse form body")
)

// Application layer (post-parse operations: validation/crypto/persist)
var (
	ErrEncrypt      = errors.New("callback: encrypt failed")
	ErrInvalidInput = errors.New("callback: invalid input")
	ErrPersist      = errors.New("callback: persist failed")
)
