package store

import (
	"errors"
)

var (
	ErrAlreadyRevoked = errors.New("already revoked")
	ErrConflict       = errors.New("conflict")
	ErrInvalidID      = errors.New("invalid refreshID")
	ErrNotFound       = errors.New("not found")
)
