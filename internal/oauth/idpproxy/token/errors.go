package token

import "errors"

var (
	ErrInvalidClient        = errors.New("token: invalid client")
	ErrInvalidGrant         = errors.New("token: invalid grant")
	ErrUnsupportedGrantType = errors.New("token: unsupported grant_type")
)
