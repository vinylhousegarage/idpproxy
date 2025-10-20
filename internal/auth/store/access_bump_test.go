package store

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_Bump(t *testing.T) {
	requireEmulator(t)
	t.Parallel()

	t.Run("empty userID -> ErrInvalidID", func(t *testing.T) {
		t.Parallel()
		r := newTestRepo(t)

		_, err := r.Bump(context.Background(), "", time.Unix(1_800_000_000, 0).UTC())
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("create new doc when absent -> gen=1 and updated_at=t", func(t *testing.T) {
		t.Parallel()
		r := newTestRepo(t)

		ctx := context.Background()
		user := "user-bump-create"
		t.Cleanup(func() { deleteAccessGenDoc(t, r, user) })

		t0 := time.Unix(1_800_000_000, 0).UTC()

		_, _ = r.docAG(user).Delete(ctx)

		gen, err := r.Bump(ctx, user, t0)
		require.NoError(t, err)
		require.Equal(t, 1, gen)

		got := getAccessGenDoc(t, r, user)
		require.Equal(t, user, got.UserID)
		require.Equal(t, 1, got.Gen)
		require.True(t, got.UpdatedAt.Equal(t0), "updated_at should equal the passed time (UTC)")
	})

	t.Run("increment existing doc -> gen increments and updated_at replaced", func(t *testing.T) {
		t.Parallel()
		r := newTestRepo(t)

		ctx := context.Background()
		user := "user-bump-increment"
		t.Cleanup(func() { deleteAccessGenDoc(t, r, user) })

		t0 := time.Unix(1_800_000_000, 0).UTC()
		gen1, err := r.Bump(ctx, user, t0)
		require.NoError(t, err)
		require.Equal(t, 1, gen1)

		t1 := t0.Add(5 * time.Minute).UTC()
		gen2, err := r.Bump(ctx, user, t1)
		require.NoError(t, err)
		require.Equal(t, 2, gen2)

		got := getAccessGenDoc(t, r, user)
		require.Equal(t, 2, got.Gen)
		require.True(t, got.UpdatedAt.Equal(t1), "updated_at should be replaced by the latest time")
	})

	t.Run("concurrent bumps -> atomic increments result in final gen == N", func(t *testing.T) {
		t.Parallel()
		r := newTestRepo(t)

		ctx := context.Background()
		user := "user-bump-concurrent"
		t.Cleanup(func() { deleteAccessGenDoc(t, r, user) })

		_, _ = r.docAG(user).Delete(ctx)

		const N = 10
		var wg sync.WaitGroup
		wg.Add(N)

		tFixed := time.Unix(1_900_000_000, 0).UTC()

		results := make([]int, N)
		errs := make([]error, N)

		for i := 0; i < N; i++ {
			i := i
			go func() {
				defer wg.Done()
				g, err := r.Bump(ctx, user, tFixed)
				results[i] = g
				errs[i] = err
			}()
		}
		wg.Wait()

		for i := 0; i < N; i++ {
			require.NoError(t, errs[i], "goroutine %d failed", i)
		}

		got := getAccessGenDoc(t, r, user)
		require.Equal(t, N, got.Gen, "final gen must equal the number of bumps")
		require.True(t, got.UpdatedAt.Equal(tFixed), "updated_at should be the last passed time (all equal here)")

		sorted := append([]int(nil), results...)
		sort.Ints(sorted)
		require.Equal(t, 1, sorted[0], "smallest newGen should be 1")
		require.Equal(t, N, sorted[len(sorted)-1], "largest newGen should be N")
	})
}
