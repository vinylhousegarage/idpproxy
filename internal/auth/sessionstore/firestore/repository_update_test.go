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
)

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("success_updates_document", func(t *testing.T) {
		t.Parallel()

		if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
			t.Skip("FIRESTORE_EMULATOR_HOST is not set (Firestore emulator required)")
		}

		ctx := context.Background()

		projectID := "idpproxy-test"
		client, err := firestore.NewClient(ctx, projectID)
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.Close() })

		collectionName := fmt.Sprintf("sessions_%d", time.Now().UnixNano())
		repo := NewRepository(client, collectionName)

		sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())

		original := &session.Session{
			SessionID: sessionID,
			UserID:    "user-1",
		}
		err = repo.Create(ctx, original)
		require.NoError(t, err)

		updated := &session.Session{
			SessionID: sessionID,
			UserID:    "user-2",
		}
		err = repo.Update(ctx, updated)
		require.NoError(t, err)

		gotDoc, err := client.Collection(collectionName).Doc(sessionID).Get(ctx)
		require.NoError(t, err)

		var got session.Session
		err = gotDoc.DataTo(&got)
		require.NoError(t, err)

		require.Equal(t, updated.SessionID, got.SessionID)
		require.Equal(t, updated.UserID, got.UserID)
	})

	t.Run("returns_error_when_ctx_canceled", func(t *testing.T) {
		t.Parallel()

		if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
			t.Skip("FIRESTORE_EMULATOR_HOST is not set (Firestore emulator required)")
		}

		parent := context.Background()
		ctx, cancel := context.WithCancel(parent)
		cancel()

		projectID := "idpproxy-test"
		client, err := firestore.NewClient(parent, projectID)
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.Close() })

		collectionName := fmt.Sprintf("sessions_%d", time.Now().UnixNano())
		repo := NewRepository(client, collectionName)

		s := &session.Session{
			SessionID: fmt.Sprintf("session-%d", time.Now().UnixNano()),
			UserID:    "user-1",
		}

		err = repo.Update(ctx, s)
		require.Error(t, err)
	})
}
