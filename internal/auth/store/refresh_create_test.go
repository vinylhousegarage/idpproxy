package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

func requireEmulator(t *testing.T) {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set; skipping Firestore emulator tests")
	}
}

func newTestRepo(t *testing.T) *Repo {
	t.Helper()
	ctx := context.Background()

	projectID := os.Getenv("TEST_FIRESTORE_PROJECT")
	if projectID == "" {
		projectID = "demo-idpproxy"
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithoutAuthentication())
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	fixed := time.Unix(1_725_000_000, 0)
	return &Repo{fs: client, now: func() time.Time { return fixed }}
}

func makeRec(id, user, fam string, now time.Time) *RefreshTokenRecord {
	return &RefreshTokenRecord{
		RefreshID:  id,
		UserID:     user,
		DigestB64:  "dummy-digest",
		KeyID:      "kid-1",
		FamilyID:   fam,
		CreatedAt:  now,
		LastUsedAt: now,
		ExpiresAt:  now.Add(24 * time.Hour),
		DeleteAt:   now.Add(48 * time.Hour),
	}
}

func TestRepo_Create(t *testing.T) {
	t.Parallel()
	requireEmulator(t)

	ctx := context.Background()
	repo := newTestRepo(t)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fixed := repo.now()
		id := fmt.Sprintf("rt-%d-success", fixed.UnixNano())
		rec := makeRec(id, "user-1", "fam-1", fixed)

		require.NoError(t, repo.Create(ctx, rec))
		t.Cleanup(func() { _, _ = repo.docRT(id).Delete(ctx) })

		got, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, rec.RefreshID, got.RefreshID)
		require.Equal(t, rec.UserID, got.UserID)
		require.Equal(t, rec.FamilyID, got.FamilyID)
	})

	t.Run("conflict", func(t *testing.T) {
		t.Parallel()

		fixed := repo.now()
		id := fmt.Sprintf("rt-%d-conflict", fixed.UnixNano())
		rec := makeRec(id, "user-2", "fam-2", fixed)

		require.NoError(t, repo.Create(ctx, rec))
		t.Cleanup(func() { _, _ = repo.docRT(id).Delete(ctx) })

		err := repo.Create(ctx, rec)
		require.ErrorIs(t, err, ErrConflict)
	})
}
