package apierror

import "errors"

var (
	ErrEncrypt      = errors.New("encrypt failed")
	ErrInvalidInput = errors.New("invalid input")
	ErrPersist      = errors.New("persist failed")
)
