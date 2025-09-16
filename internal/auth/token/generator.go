package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/auth/store"
)

const (
	refreshTokenPrefix = "rt1."
	refreshTokenRawLen = 32
)

func GenerateRefreshToken(
	_ context.Context,
	userID string,
	ttl time.Duration,
	purgeAfter time.Duration,
) (*store.RefreshTokenRecord, string, error) {
	if userID == "" {
		return nil, "", errors.New("userID empty")
	}
	if ttl <= 0 {
		return nil, "", errors.New("ttl must be > 0")
	}
	if purgeAfter < ttl {
		return nil, "", errors.New("purgeAfter must be >= ttl")
	}

	raw := make([]byte, refreshTokenRawLen)
	if _, err := rand.Read(raw); err != nil {
		return nil, "", fmt.Errorf("rand: %w", err)
	}

	now := time.Now().UTC()
	rec := &store.RefreshTokenRecord{
		UserID:     userID,
		CreatedAt:  now,
		LastUsedAt: now,
		ExpiresAt:  now.Add(ttl),
		DeleteAt:   now.Add(purgeAfter),
	}

	token := refreshTokenPrefix + base64.RawURLEncoding.EncodeToString(raw)
	return rec, token, nil
}
