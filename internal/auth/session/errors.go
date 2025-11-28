package session

import "errors"

// Usecase validation
var (
	ErrEmptySessionID       = errors.New("session: empty sessionID")
	ErrEmptyUserID          = errors.New("session: empty userID")
	ErrInvalidUsecaseConfig = errors.New("session: invalid usecase configuration")
)

// Repository layer
var (
	ErrNotFound = errors.New("session: not found")
)
