package firesessionstore

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"google.golang.org/api/option"
)

func TestRepository_PurgeExpired(t *testing.T) {
	t.Parallel()

	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set (Firestore emulator is required)")
	}

	t.Run("purge_expired_only_and_return_deleted_count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		client := newFirestoreClient(t, ctx)
		defer client.Close()

		collection := uniqueCollectionName("sessions")
		repo := NewRepository(client, collection)

		now := time.Now().UTC()
		before := now

		expired1 := &session.Session{
			SessionID: "session-expired-1",
			UserID:    "user-1",
			Status:    "active",
			ExpiresAt: now.Add(-2 * time.Hour),
			CreatedAt: now.Add(-3 * time.Hour),
		}
		expired2 := &session.Session{
			SessionID: "session-expired-2",
			UserID:    "user-2",
			Status:    "active",
			ExpiresAt: now.Add(-1 * time.Minute),
			CreatedAt: now.Add(-10 * time.Minute),
		}
		active := &session.Session{
			SessionID: "session-active-1",
			UserID:    "user-3",
			Status:    "active",
			ExpiresAt: now.Add(+2 * time.Hour),
			CreatedAt: now.Add(-1 * time.Minute),
		}

		require.NoError(t, repo.Create(ctx, expired1))
		require.NoError(t, repo.Create(ctx, expired2))
		require.NoError(t, repo.Create(ctx, active))

		deleted, err := repo.PurgeExpired(ctx, before)
		require.NoError(t, err)
		require.Equal(t, 2, deleted)

		_, err = repo.FindByID(ctx, expired1.SessionID)
		require.ErrorIs(t, err, session.ErrNotFound)

		_, err = repo.FindByID(ctx, expired2.SessionID)
		require.ErrorIs(t, err, session.ErrNotFound)

		got, err := repo.FindByID(ctx, active.SessionID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, active.SessionID, got.SessionID)
	})

	t.Run("no_expired_returns_zero_and_keeps_all", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		client := newFirestoreClient(t, ctx)
		defer client.Close()

		collection := uniqueCollectionName("sessions")
		repo := NewRepository(client, collection)

		now := time.Now().UTC()
		before := now

		active1 := &session.Session{
			SessionID: "session-active-1",
			UserID:    "user-1",
			Status:    "active",
			ExpiresAt: now.Add(+10 * time.Minute),
			CreatedAt: now.Add(-1 * time.Minute),
		}
		active2 := &session.Session{
			SessionID: "session-active-2",
			UserID:    "user-2",
			Status:    "active",
			ExpiresAt: now.Add(+2 * time.Hour),
			CreatedAt: now.Add(-1 * time.Minute),
		}

		require.NoError(t, repo.Create(ctx, active1))
		require.NoError(t, repo.Create(ctx, active2))

		deleted, err := repo.PurgeExpired(ctx, before)
		require.NoError(t, err)
		require.Equal(t, 0, deleted)

		_, err = repo.FindByID(ctx, active1.SessionID)
		require.NoError(t, err)
		_, err = repo.FindByID(ctx, active2.SessionID)
		require.NoError(t, err)
	})

	t.Run("context_canceled_returns_error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		client := newFirestoreClient(t, ctx)
		defer client.Close()

		collection := uniqueCollectionName("sessions")
		repo := NewRepository(client, collection)

		now := time.Now().UTC()
		expired := &session.Session{
			SessionID: "session-expired-1",
			UserID:    "user-1",
			Status:    "active",
			ExpiresAt: now.Add(-1 * time.Minute),
			CreatedAt: now.Add(-2 * time.Minute),
		}
		require.NoError(t, repo.Create(ctx, expired))

		cctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := repo.PurgeExpired(cctx, now)
		require.Error(t, err)
	})
}

func newFirestoreClient(t *testing.T, ctx context.Context) *firestore.Client {
	t.Helper()

	projectID := "test-project"

	client, err := firestore.NewClient(ctx, projectID, option.WithoutAuthentication())
	require.NoError(t, err)

	return client
}

func uniqueCollectionName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
