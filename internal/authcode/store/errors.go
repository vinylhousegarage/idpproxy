package store

import "errors"

var (
	ErrClientMismatch = errors.New("proxycode client mismatch")
	ErrExpired        = errors.New("proxycode expired")
	ErrNotFound       = errors.New("proxycode not found")
)
