package store

import (
	"errors"
)

var (
	ErrAlreadyRevoked  = errors.New("already revoked")
	ErrConflict        = errors.New("conflict")
	ErrDeleted         = errors.New("refresh token deleted")
	ErrExpired         = errors.New("refresh token expired")
	ErrInvalid         = errors.New("invalid record")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidID       = errors.New("invalid refreshID")
	ErrInvalidUntil    = errors.New("invalid until")
	ErrInvalidUserID   = errors.New("invalid UserID")
	ErrNotFound        = errors.New("not found")
	ErrRevoked         = errors.New("refresh token revoked")
)
