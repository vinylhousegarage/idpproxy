package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_Get(t *testing.T) {
	t.Parallel()
	requireEmulator(t)

	ctx := context.Background()
	r := newTestRepo(t)

	t.Run("invalid userID (empty) -> ErrInvalidID", func(t *testing.T) {
		t.Parallel()

		rec, err := r.Get(ctx, "")
		require.Nil(t, rec)
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("not found -> ErrNotFound", func(t *testing.T) {
		t.Parallel()

		const uid = "user-not-found-get-1"
		_, _ = r.fs.Collection(colAccessGenerations).Doc(uid).Delete(ctx)

		rec, err := r.Get(ctx, uid)
		require.Nil(t, rec)
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("success -> returns record", func(t *testing.T) {
		t.Parallel()

		const uid = "user-success-get-1"
		doc := r.fs.Collection(colAccessGenerations).Doc(uid)
		t.Cleanup(func() { _, _ = doc.Delete(ctx) })

		seed := AccessGenerationRecord{
			UserID:    uid,
			Gen:       3,
			UpdatedAt: time.Unix(1_900_000_000, 0).UTC(),
		}
		_, err := doc.Set(ctx, seed)
		require.NoError(t, err)

		got, err := r.Get(ctx, uid)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, uid, got.UserID)
		require.Equal(t, 3, got.Gen)
		require.WithinDuration(t, seed.UpdatedAt, got.UpdatedAt, time.Second)
	})

	t.Run("doc exists but missing user_id field -> fills from docID", func(t *testing.T) {
		t.Parallel()

		const uid = "user-missing-userid-get-1"
		doc := r.fs.Collection(colAccessGenerations).Doc(uid)
		t.Cleanup(func() { _, _ = doc.Delete(ctx) })

		_, err := doc.Set(ctx, map[string]any{
			"gen":        7,
			"updated_at": time.Unix(2_000_000_000, 0).UTC(),
		})
		require.NoError(t, err)

		got, err := r.Get(ctx, uid)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, uid, got.UserID)
		require.Equal(t, 7, got.Gen)
	})
}
