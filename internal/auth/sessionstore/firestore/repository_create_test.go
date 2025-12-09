package firesessionstore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		client, err := firestore.NewClient(ctx, "test-project")
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = client.Close()
		})

		const collectionName = "sessions_test"
		repo := NewRepository(client, collectionName)

		s := &session.Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
		}

		err = repo.Create(ctx, s)
		require.NoError(t, err)

		snap, err := client.Collection(collectionName).Doc("session-123").Get(ctx)
		require.NoError(t, err)

		var got session.Session
		err = snap.DataTo(&got)
		require.NoError(t, err)

		require.Equal(t, s.SessionID, got.SessionID)
		require.Equal(t, s.UserID, got.UserID)
		require.Equal(t, s.Status, got.Status)
	})
}
