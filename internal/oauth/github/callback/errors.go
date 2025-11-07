package callback

import "errors"

var (
	ErrInvalidInput = errors.New("github: invalid input")
	ErrEncrypt      = errors.New("github: encrypt failed")
	ErrPersist      = errors.New("github: persist failed")
)
