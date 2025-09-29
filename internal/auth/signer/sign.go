package signer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *HMACSigner) validateKey() error {
	if len(s.key) == 0 {
		return ErrEmptyKey
	}

	return nil
}

func buildClaims(payload []byte, now time.Time) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	}

	if len(payload) > 0 {
		var m map[string]any
		if err := json.Unmarshal(payload, &m); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidPayload, err)
		}
		for k, v := range m {
			claims[k] = v
		}
	}

	return claims, nil
}

func signToken(claims jwt.Claims, key []byte, kid string) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok.Header["typ"] = "JWT"
	if kid != "" {
		tok.Header["kid"] = kid
	}

	return tok.SignedString(key)
}

func (s *HMACSigner) Sign(ctx context.Context, payload []byte) (string, string, error) {
	_ = ctx

	if err := s.validateKey(); err != nil {
		return "", "", err
	}

	claims, err := buildClaims(payload, s.Now())
	if err != nil {
		return "", "", err
	}

	token, err := signToken(claims, s.key, s.keyID)
	if err != nil {
		return "", "", fmt.Errorf("sign jwt: %w", err)
	}

	return token, s.keyID, nil
}
