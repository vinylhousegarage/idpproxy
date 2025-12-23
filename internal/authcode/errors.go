package authcode

import "errors"

var (
	ErrAlreadyUsed    = errors.New("authcode already used")
	ErrClientMismatch = errors.New("client mismatch")
	ErrExpired        = errors.New("authcode expired")
	ErrNotFound       = errors.New("authcode not found")
)
