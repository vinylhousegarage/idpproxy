package token

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/auth/store"
)

var timeNow = func() time.Time { return time.Now().UTC() }

func currentPepperKeyID() (string, error) {
	v := os.Getenv("IDPPROXY_REFRESH_PEPPER_KEY_ID")
	if v == "" {
		return "", fmt.Errorf("IDPPROXY_REFRESH_PEPPER_KEY_ID is required but not set")
	}

	return v, nil
}

var getPepperKeyMaterial = func(keyID string) ([]byte, error) {
	v := os.Getenv("IDPPROXY_REFRESH_PEPPER_KEY_MATERIAL")
	if v == "" {
		return nil, fmt.Errorf("IDPPROXY_REFRESH_PEPPER_KEY_MATERIAL is required but not set (KeyID=%s)", keyID)
	}

	return []byte(v), nil
}

func computeDigestB64(secretRaw []byte, keyID string) string {
	key := getPepperKeyMaterial(keyID)
	mac := hmac.New(sha256.New, key)
	mac.Write(secretRaw)
	sum := mac.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(sum)
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

	familyID := newFamilyID()

	secretRaw := make([]byte, refreshTokenRawLen)
	if _, err := rand.Read(secretRaw); err != nil {
		return nil, "", errors.Join(ErrRandFailure, err)
	}

	keyID, err := currentPepperKeyID()
	if err != nil {
		return nil, "", fmt.Errorf("pepper key id: %w", err)
	}

	digestB64, err := computeDigestB64(secretRaw, keyID)
	if err != nil {
		return nil, "", fmt.Errorf("compute digest: %w", err)
	}

	secretB64 := base64.RawURLEncoding.EncodeToString(secretRaw)
	token := fmt.Sprintf("%s%s.%s", refreshTokenPrefix, refreshID, secretB64)

	now := timeNow()
	rec := &store.RefreshTokenRecord{
		RefreshID: refreshID,
		UserID:    userID,
		DigestB64: digestB64,
		KeyID:     keyID,

		FamilyID: familyID,

		CreatedAt:  now,
		LastUsedAt: time.Time{},
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
