package token

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testUserID1 = "user-1"
	testUserID2 = "user-2"
)

func TestGenerateRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("basic properties (prefix, split, lengths, digest, family, times)", func(t *testing.T) {
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
		require.Equal(t, idB64, rec.RefreshID)

		sum := sha256.Sum256(secretRaw)
		expectDigest := base64.RawURLEncoding.EncodeToString(sum[:])
		require.Equal(t, expectDigest, rec.DigestB64, "DigestB64 should be sha256(secret) base64url")

		require.Equal(t, rec.RefreshID, rec.FamilyID, "initial token should set FamilyID = RefreshID")

		require.WithinDuration(t, rec.CreatedAt.Add(ttl), rec.ExpiresAt, time.Second)
		require.WithinDuration(t, rec.CreatedAt.Add(purge), rec.DeleteAt, time.Second)
		require.WithinDuration(t, rec.CreatedAt, rec.LastUsedAt, time.Second)
	})

	t.Run("randomness (tokens and IDs differ)", func(t *testing.T) {
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
	})
}
