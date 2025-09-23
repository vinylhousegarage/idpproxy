package jwt

import "slices"

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
