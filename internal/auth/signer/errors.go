package signer

import "errors"

// Sign + Verify
var (
	ErrEmptyKey = errors.New("hmacsigner: empty key")
)

// Sign
var (
	ErrInvalidPayload = errors.New("hmacsigner: invalid payload json")
)

// Verify
var (
	ErrEmptyToken    = errors.New("hmacsigner: empty token")
	ErrInvalidAlg    = errors.New("hmacsigner: invalid alg")
	ErrInvalidToken  = errors.New("hmacsigner: invalid token")
	ErrInvalidTyp    = errors.New("hmacsigner: invalid typ")
	ErrUnexpectedKID = errors.New("hmacsigner: unexpected kid")
)
