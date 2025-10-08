package token

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testUserID1 = "github:12345678"
	testUserID2 = "github:87654321"
)

func TestGenerateRefreshToken(t *testing.T) {
	t.Parallel()

	oldKeyID := currentPepperKeyID
	oldKeyMat := getPepperKeyMaterial
	oldNow := timeNow

	const testKeyID = "test-hmac-k1"
	const testPepper = "pepper-for-test"
	fixedNow := time.Unix(1_800_000_000, 0).UTC()

	currentPepperKeyID = func() string { return testKeyID }
	getPepperKeyMaterial = func(keyID string) []byte { return []byte(testPepper) }
	timeNow = func() time.Time { return fixedNow }

	t.Cleanup(func() {
		currentPepperKeyID = oldKeyID
		getPepperKeyMaterial = oldKeyMat
		timeNow = oldNow
	})

	t.Run("basic properties (format, lengths, digest=HMAC, keyID, family=uuid, times, lastUsed=zero)", func(t *testing.T) {
		t.Parallel()

		ttl := 24 * time.Hour
		purge := 48 * time.Hour

		rec, token, err := GenerateRefreshToken(context.Background(), testUserID1, ttl, purge)
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(token, refreshTokenPrefix))
		rest := strings.TrimPrefix(token, refreshTokenPrefix)
		parts := strings.Split(rest, ".")
		require.Len(t, parts, 2, "token should be '<prefix><id>.<secretB64>'")

		idB64 := parts[0]
		secretB64 := parts[1]

		idRaw, err := base64.RawURLEncoding.DecodeString(idB64)
		require.NoError(t, err)
		require.Len(t, idRaw, refreshIDRawLen)

		secretRaw, err := base64.RawURLEncoding.DecodeString(secretB64)
		require.NoError(t, err)
		require.Len(t, secretRaw, refreshTokenRawLen)

		require.Equal(t, testUserID1, rec.UserID)
		require.Equal(t, idB64, rec.RefreshID, "RefreshID should equal idB64 part")
		require.NotEmpty(t, rec.KeyID)
		require.Equal(t, testKeyID, rec.KeyID)

		mac := hmac.New(sha256.New, []byte(testPepper))
		mac.Write(secretRaw)
		expectDigest := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		require.Equal(t, expectDigest, rec.DigestB64, "DigestB64 should be HMAC-SHA256(pepper, secret) base64url")

		require.NotEmpty(t, rec.FamilyID)
		_, err = uuid.Parse(rec.FamilyID)
		require.NoError(t, err, "FamilyID should be a UUID")

		require.WithinDuration(t, fixedNow, rec.CreatedAt, time.Second)
		require.True(t, rec.LastUsedAt.IsZero(), "LastUsedAt should be zero at issuance")
		require.WithinDuration(t, rec.CreatedAt.Add(ttl), rec.ExpiresAt, time.Second)
		require.WithinDuration(t, rec.CreatedAt.Add(purge), rec.DeleteAt, time.Second)
	})

	t.Run("randomness (opaque token, id, digest differ per issuance)", func(t *testing.T) {
		t.Parallel()

		ttl := time.Hour
		purge := 2 * time.Hour

		rec1, tok1, err1 := GenerateRefreshToken(context.Background(), testUserID2, ttl, purge)
		require.NoError(t, err1)

		rec2, tok2, err2 := GenerateRefreshToken(context.Background(), testUserID2, ttl, purge)
		require.NoError(t, err2)

		require.NotEqual(t, tok1, tok2, "opaque tokens should differ")
		require.NotEqual(t, rec1.RefreshID, rec2.RefreshID, "RefreshID should be random per issuance")
		require.NotEqual(t, rec1.DigestB64, rec2.DigestB64, "Digest should differ because secret differs")
		require.Equal(t, rec1.KeyID, rec2.KeyID, "KeyID should be same current pepper key in test")
	})

	t.Run("validation errors (empty user, invalid ttl/purge)", func(t *testing.T) {
		t.Parallel()

		_, _, err := GenerateRefreshToken(context.Background(), "", time.Hour, 2*time.Hour)
		require.Error(t, err, "empty userID should error")

		_, _, err = GenerateRefreshToken(context.Background(), testUserID1, 0, time.Hour)
		require.Error(t, err, "ttl <= 0 should error")

		_, _, err = GenerateRefreshToken(context.Background(), testUserID1, time.Hour, 30*time.Minute)
		require.Error(t, err, "purgeAfter < ttl should error")
	})
}
