package callback

import "errors"

var (
	ErrEncrypt      = errors.New("callback: encrypt failed")
	ErrInvalidInput = errors.New("callback: invalid input")
	ErrPersist      = errors.New("callback: persist failed")
)
