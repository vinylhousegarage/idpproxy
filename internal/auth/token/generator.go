package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
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

	idRaw := make([]byte, refreshIDRawLen)
	if _, err := rand.Read(idRaw); err != nil {
		return nil, "", errors.Join(ErrRandFailure, err)
	}
	refreshID := base64.RawURLEncoding.EncodeToString(idRaw)

	secretRaw := make([]byte, refreshTokenRawLen)
	if _, err := rand.Read(secretRaw); err != nil {
		return nil, "", errors.Join(ErrRandFailure, err)
	}

	sum := sha256.Sum256(secretRaw)
	digestB64 := base64.RawURLEncoding.EncodeToString(sum[:])

	secretB64 := base64.RawURLEncoding.EncodeToString(secretRaw)
	token := fmt.Sprintf("%s%s.%s", refreshTokenPrefix, refreshID, secretB64)

	now := timeNow()
	rec := &store.RefreshTokenRecord{
		RefreshID: refreshID,
		UserID:    userID,
		DigestB64: digestB64,
		KeyID:     "",
		FamilyID:  refreshID,

		CreatedAt:  now,
		LastUsedAt: now,
		ExpiresAt:  now.Add(ttl),
		DeleteAt:   now.Add(purgeAfter),
	}

	for i := range secretRaw {
		secretRaw[i] = 0
	}
	for i := range idRaw {
		idRaw[i] = 0
	}

	return rec, token, nil
}
