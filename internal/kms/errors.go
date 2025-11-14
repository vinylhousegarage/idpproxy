package kms

import "errors"

// Adapter-level errors
var (
	ErrBadFormat     = errors.New("kms: bad ciphertext format")
	ErrDecryptFailed = errors.New("kms: decrypt failed")
	ErrEncryptFailed = errors.New("kms: encrypt failed")
)

// Client-level errors
var (
	ErrInitFailed = errors.New("kms: init failed")
)
