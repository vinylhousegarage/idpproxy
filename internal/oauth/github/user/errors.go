package user

import "errors"

var (
	ErrEmptyBearerToken                 = errors.New("empty bearer token")
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid Authorization header format")
	ErrMissingAuthorizationHeader       = errors.New("missing Authorization header")
	ErrNilContext                       = errors.New("nil context")
)
