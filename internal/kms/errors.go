package kms

import "errors"

var (
	ErrBadFormat     = errors.New("kms: bad ciphertext format")
	ErrDecryptFailed = errors.New("kms: decrypt failed")
	ErrEncryptFailed = errors.New("kms: encrypt failed")
)
