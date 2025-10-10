package store

import (
	"errors"
)

var (
	ErrAlreadyRevoked = errors.New("already revoked")
	ErrConflict       = errors.New("conflict")
	ErrDeleted        = errors.New("refresh token deleted")
	ErrExpired        = errors.New("refresh token expired")
	ErrInvalidID      = errors.New("invalid refreshID")
	ErrInvalidUserID  = errors.New("invalid UserID")
	ErrNotFound       = errors.New("not found")
	ErrRevoked        = errors.New("refresh token revoked")
)
