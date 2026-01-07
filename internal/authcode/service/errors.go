package service

import "errors"

var (
	ErrClientMismatch = errors.New("authcode client mismatch")
	ErrExpired        = errors.New("authcode expired")
)
