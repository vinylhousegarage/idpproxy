package idtoken

import (
	"errors"
)

var (
	ErrInvalidAudience = errors.New("invalid issuer")
	ErrInvalidExp      = errors.New("invalid expiration")
	ErrInvalidIat      = errors.New("invalid issued-at")
	ErrInvalidIssuer   = errors.New("invalid issuer")
	ErrInvalidSubject  = errors.New("invalid subject")
)
