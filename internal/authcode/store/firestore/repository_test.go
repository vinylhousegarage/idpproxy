package fireauthcodestore

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
	authcodestore "github.com/vinylhousegarage/idpproxy/internal/authcode/store"
)

func TestRepository_Consume(t *testing.T) {
	ctx := context.Background()
	now := time.Date(
		2026,
		time.November,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	t.Run("returns proxy code and deletes it after consumption", func(t *testing.T) {
		repository := newTestRepository(t, now)

		want := authcode.ProxyCode{
			Code:      "proxy-code-123",
			UserID:    "user-123",
			ClientID:  "client-123",
			ExpiresAt: now.Add(time.Minute),
		}
		require.NoError(t, repository.Save(ctx, want))

		got, err := repository.Consume(
			ctx,
			want.Code,
			want.ClientID,
		)

		require.NoError(t, err)
		require.Equal(t, want, got)

		_, err = repository.Consume(
			ctx,
			want.Code,
			want.ClientID,
		)
		require.ErrorIs(t, err, authcodestore.ErrNotFound)
	})

	t.Run("returns not found when proxy code does not exist", func(t *testing.T) {
		repository := newTestRepository(t, now)

		_, err := repository.Consume(
			ctx,
			"missing-proxy-code",
			"client-123",
		)

		require.ErrorIs(t, err, authcodestore.ErrNotFound)
	})

	t.Run("returns client mismatch and preserves proxy code", func(t *testing.T) {
		repository := newTestRepository(t, now)

		proxyCode := authcode.ProxyCode{
			Code:      "proxy-code-123",
			UserID:    "user-123",
			ClientID:  "client-123",
			ExpiresAt: now.Add(time.Minute),
		}
		require.NoError(t, repository.Save(ctx, proxyCode))

		_, err := repository.Consume(
			ctx,
			proxyCode.Code,
			"other-client",
		)
		require.ErrorIs(t, err, authcodestore.ErrClientMismatch)

		got, err := repository.Consume(
			ctx,
			proxyCode.Code,
			proxyCode.ClientID,
		)
		require.NoError(t, err)
		require.Equal(t, proxyCode, got)
	})

	t.Run("returns expired and deletes an expired proxy code", func(t *testing.T) {
		repository := newTestRepository(t, now)

		proxyCode := authcode.ProxyCode{
			Code:      "expired-proxy-code",
			UserID:    "user-123",
			ClientID:  "client-123",
			ExpiresAt: now,
		}
		require.NoError(t, repository.Save(ctx, proxyCode))

		_, err := repository.Consume(
			ctx,
			proxyCode.Code,
			proxyCode.ClientID,
		)
		require.ErrorIs(t, err, authcodestore.ErrExpired)

		_, err = repository.Consume(
			ctx,
			proxyCode.Code,
			proxyCode.ClientID,
		)
		require.ErrorIs(t, err, authcodestore.ErrNotFound)
	})
}

func newTestRepository(
	t *testing.T,
	now time.Time,
) *Repository {
	t.Helper()

	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	client, err := firestore.NewClient(
		context.Background(),
		"idpproxy-test",
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, client.Close())
	})

	repository := NewRepository(
		client,
		fmt.Sprintf(
			"proxy_codes_test_%d",
			time.Now().UnixNano(),
		),
	)
	repository.now = func() time.Time {
		return now
	}

	return repository
}
