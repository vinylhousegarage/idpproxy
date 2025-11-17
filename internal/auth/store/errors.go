package store

import (
	"errors"
)

// Access token + Refresh token
var (
	ErrInvalid   = errors.New("token record invalid")
	ErrInvalidID = errors.New("token id invalid")
	ErrNotFound  = errors.New("token record not found")
)

// Access token
var (
	ErrInvalidArgument = errors.New("access token invalid argument")
)

// Refresh token
var (
	ErrAlreadyRevoked = errors.New("refresh token already revoked")
	ErrConflict       = errors.New("refresh token conflict")
	ErrDeleted        = errors.New("refresh token deleted")
	ErrExpired        = errors.New("refresh token expired")
	ErrInvalidUntil   = errors.New("refresh token invalid until")
	ErrRevoked        = errors.New("refresh token revoked")
)

// Validation
var (
	ErrInvalidUserID = errors.New("invalid user id")
)
