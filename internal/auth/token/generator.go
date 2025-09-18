package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/auth/store"
)

var timeNow = func() time.Time { return time.Now().UTC() }

func validateParams(userID string, ttl, purgeAfter time.Duration) error {
	switch {
	case userID == "":
		return ErrEmptyUserID
	case ttl <= 0:
		return ErrInvalidTTL
	case purgeAfter < ttl:
		return ErrInvalidPurge
	default:
		return nil
	}
}

func GenerateRefreshToken(ctx context.Context, userID string, ttl, purgeAfter time.Duration) (*store.RefreshTokenRecord, string, error) {
	_ = ctx

	if err := validateParams(userID, ttl, purgeAfter); err != nil {
		return nil, "", err
	}

	raw := make([]byte, refreshTokenRawLen)
	if _, err := rand.Read(raw); err != nil {
		return nil, "", errors.Join(ErrRandFailure, err)
	}
	defer func() {
		for i := range raw {
			raw[i] = 0
		}
	}()

	now := timeNow()
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
