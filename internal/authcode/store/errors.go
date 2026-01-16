package store

import "errors"

var (
	ErrClientMismatch = errors.New("authcode client mismatch")
	ErrExpired        = errors.New("authcode expired")
	ErrNotFound       = errors.New("authcode not found")
)
