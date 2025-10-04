package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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

		snap, err := repo.docRT(id).Get(ctx)
		require.NoError(t, err)

		var got RefreshTokenRecord
		require.NoError(t, snap.DataTo(&got))
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
