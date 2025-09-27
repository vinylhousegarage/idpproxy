package signer

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *HMACSigner) Sign(ctx context.Context, payload []byte) (string, string, error) {
	_ = ctx

	if len(s.key) == 0 {
		return "", "", ErrEmptyKey
	}

	now := s.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	}

	if len(payload) > 0 {
		var m map[string]any
		if err := json.Unmarshal(payload, &m); err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrInvalidPayload, err)
		}
		for k, v := range m {
			claims[k] = v
		}
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok.Header["typ"] = "JWT"
	if s.keyID != "" {
		tok.Header["kid"] = s.keyID
	}

	signed, err := tok.SignedString(s.key)
	if err != nil {
		return "", "", fmt.Errorf("sign jwt: %w", err)
	}

	return signed, s.keyID, nil
}
