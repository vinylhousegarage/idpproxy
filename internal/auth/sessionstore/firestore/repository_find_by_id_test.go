package firesessionstore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

func TestRepository_FindByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "test-project")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = client.Close()
	})

	collectionName := fmt.Sprintf("sessions_findbyid_test_%d", time.Now().UnixNano())
	repo := NewRepository(client, collectionName)

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		s := &session.Session{
			SessionID: "session-123",
			UserID:    "user-123",
		}

		err := repo.Create(ctx, s)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "session-123")
		require.NoError(t, err)
		require.NotNil(t, got)

		require.Equal(t, s.SessionID, got.SessionID)
		require.Equal(t, s.UserID, got.UserID)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		got, err := repo.FindByID(ctx, "no-such-session")
		require.Nil(t, got)
		require.ErrorIs(t, err, session.ErrNotFound)
	})
}
