package encrypt

import "errors"

var (
	ErrBadFormat     = errors.New("ciphertext format invalid")
	ErrDecryptFailed = errors.New("decrypt failed")
	ErrEmptyBlob     = errors.New("empty ciphertext blob")
	ErrNilKey        = errors.New("nil key")
	ErrNilKeySet     = errors.New("nil or empty key set")
	ErrUnknownKID    = errors.New("unknown KID")
)
