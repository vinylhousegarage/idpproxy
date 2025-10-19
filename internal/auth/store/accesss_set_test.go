package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func getAccessGenDoc(t *testing.T, r *Repo, user string) *AccessGenerationRecord {
	t.Helper()
	ctx := context.Background()

	snap, err := r.docAG(user).Get(ctx)
	require.NoError(t, err)

	var got AccessGenerationRecord
	require.NoError(t, snap.DataTo(&got))
	return &got
}

func TestRepo_Set(t *testing.T) {
	requireEmulator(t)
	t.Parallel()

	t.Run("nil record -> ErrInvalidArgument", func(t *testing.T) {
		t.Parallel()
		r := newTestRepo(t)

		err := r.Set(context.Background(), nil)
		require.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("empty UserID -> ErrInvalidID", func(t *testing.T) {
		t.Parallel()
		r := newTestRepo(t)

		err := r.Set(context.Background(), &AccessGenerationRecord{})
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("create new doc (upsert) and set UpdatedAt to now()", func(t *testing.T) {
		t.Parallel()

		fixed := time.Unix(1_900_000_000, 0).UTC()
		r := newTestRepoWithNow(t, fixed)

		user := "u-set-create-" + t.Name()
		ctx := context.Background()
		t.Cleanup(func() { _, _ = r.docAG(user).Delete(ctx) })

		rec := &AccessGenerationRecord{UserID: user}

		require.NoError(t, r.Set(ctx, rec))

		got := getAccessGenDoc(t, r, user)
		require.Equal(t, user, got.UserID)
		require.True(t, got.UpdatedAt.Equal(fixed), "UpdatedAt should be set to repo.now()")
	})

	t.Run("update existing doc (merge) and bump UpdatedAt", func(t *testing.T) {
		t.Parallel()

		old := time.Unix(1_850_000_000, 0).UTC()
		r1 := newTestRepoWithNow(t, old)

		user := "u-set-update-" + t.Name()
		ctx := context.Background()
		t.Cleanup(func() { _, _ = r1.docAG(user).Delete(ctx) })

		require.NoError(t, r1.Set(ctx, &AccessGenerationRecord{UserID: user}))

		got1 := getAccessGenDoc(t, r1, user)
		require.True(t, got1.UpdatedAt.Equal(old))

		newer := time.Unix(1_860_000_000, 0).UTC()
		r2 := newTestRepoWithNow(t, newer)

		require.NoError(t, r2.Set(ctx, &AccessGenerationRecord{UserID: user}))

		got2 := getAccessGenDoc(t, r2, user)
		require.True(t, got2.UpdatedAt.After(got1.UpdatedAt), "UpdatedAt should advance on update")
		require.True(t, got2.UpdatedAt.Equal(newer), "UpdatedAt should be the new repo.now()")
	})
}
