package signer

import "slices"

const AlgHS256 = "HS256"

type HMACSigner struct {
	key   []byte
	keyID string
}

func NewHMACSigner(key []byte, keyID string) *HMACSigner {
	keyCopy := slices.Clone(key)
	return &HMACSigner{
		key:   keyCopy,
		keyID: keyID,
	}
}

func (s *HMACSigner) Alg() string   { return AlgHS256 }
func (s *HMACSigner) KeyID() string { return s.keyID }
