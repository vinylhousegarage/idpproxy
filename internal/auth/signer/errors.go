package signer

import "errors"

var (
	ErrEmptyKey       = errors.New("hmac signer: empty key")
	ErrInvalidPayload = errors.New("hmac signer: invalid payload json")
)
