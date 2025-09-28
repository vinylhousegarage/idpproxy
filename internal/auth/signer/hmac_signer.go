package signer

import (
	"slices"
	"time"
)

const AlgHS256 = "HS256"

type HMACSigner struct {
	key   []byte
	keyID string
	now   func() time.Time
}

func NewHMACSigner(key []byte, keyID string) *HMACSigner {
	return &HMACSigner{
		key:   slices.Clone(key),
		keyID: keyID,
		now:   time.Now,
	}
}

func (s *HMACSigner) Alg() string { return AlgHS256 }

func (s *HMACSigner) KeyID() string { return s.keyID }

func (s *HMACSigner) Now() time.Time {
	if s.now != nil {
		return s.now()
	}

	return time.Now()
}
