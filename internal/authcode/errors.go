package authcode

import "errors"

var (
	ErrAlreadyUsed = errors.New("authcode already used")
)
