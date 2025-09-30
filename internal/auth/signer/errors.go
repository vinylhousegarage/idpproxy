package signer

import "errors"

var (
	ErrEmptyKey       = errors.New("hmacsigner: empty key")
	ErrEmptyToken     = errors.New("hmacsigner: empty token")
	ErrInvalidAlg     = errors.New("hmacsigner: invalid alg")
	ErrInvalidPayload = errors.New("hmacsigner: invalid payload json")
	ErrInvalidToken   = errors.New("hmacsigner: invalid token")
	ErrInvalidTyp     = errors.New("hmacsigner: invalid typ")
	ErrUnexpectedKID  = errors.New("hmacsigner: unexpected kid")
)
