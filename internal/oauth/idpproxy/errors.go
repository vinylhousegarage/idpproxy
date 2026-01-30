package token

import "errors"

var (
	ErrExpiredCode          = errors.New("token: expired auth code")
	ErrInvalidClient        = errors.New("token: invalid client")
	ErrInvalidGrant         = errors.New("token: invalid grant")
	ErrUnsupportedGrantType = errors.New("token: unsupported grant_type")
)
