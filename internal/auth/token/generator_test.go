package token

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	refreshTokenB64Len = 43
)

func TestGenerateRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("basic properties", func(t *testing.T) {
		t.Parallel()

		ttl := 24 * time.Hour
		purge := 48 * time.Hour

		rec, token, err := GenerateRefreshToken(context.Background(), "user-1", ttl, purge)
		require.NoError(t, err)

		require.True(t, strings.HasPrefix(token, refreshTokenPrefix))
		rawB64 := strings.TrimPrefix(token, refreshTokenPrefix)

		require.Len(t, rawB64, tokenB64Len)

		raw, err := base64.RawURLEncoding.DecodeString(rawB64)
		require.NoError(t, err)
		require.Len(t, raw, refreshTokenRawLen)

		require.Equal(t, "user-1", rec.UserID)
		require.WithinDuration(t, rec.CreatedAt.Add(ttl), rec.ExpiresAt, time.Second)
		require.WithinDuration(t, rec.CreatedAt.Add(purge), rec.DeleteAt, time.Second)
		require.WithinDuration(t, rec.CreatedAt, rec.LastUsedAt, time.Second)
	})

	t.Run("randomness", func(t *testing.T) {
		t.Parallel()

		ttl := time.Hour
		purge := 2 * time.Hour

		_, tok1, err1 := GenerateRefreshToken(context.Background(), "u", ttl, purge)
		require.NoError(t, err1)

		_, tok2, err2 := GenerateRefreshToken(context.Background(), "u", ttl, purge)
		require.NoError(t, err2)

		require.NotEqual(t, tok1, tok2)
	})
}
