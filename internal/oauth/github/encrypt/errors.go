package encrypt

import "errors"

var (
	ErrNilKey     = errors.New("fernet key is nil")
	ErrEmptyToken = errors.New("token is empty")
)
