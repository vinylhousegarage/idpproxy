package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_RevokeFamily(t *testing.T) {
	requireEmulator(t)
	t.Parallel()

	fixed := time.Unix(1_800_000_000, 0).UTC()
	repo := newTestRepoWithNow(t, fixed)
	ctx := context.Background()

	t.Run("success: revoke only target family and only non-revoked", func(t *testing.T) {
		t.Parallel()

		const famA = "fam-A"
		const famB = "fam-B"

		a1 := makeActiveRec("rt-a1", "github:11111111-1111-1111-1111-111111111111", fixed)
		a1.FamilyID = famA

		a2 := makeActiveRec("rt-a2", "github:22222222-2222-2222-2222-222222222222", fixed)
		a2.FamilyID = famA
		a2.ReplacedBy = "rt-a2-new"

		a3 := makeActiveRec("rt-a3", "github:33333333-3333-3333-3333-333333333333", fixed)
		a3.FamilyID = famA
		a3.RevokedAt = fixed.Add(-time.Hour)
		a3.RevokeReason = "already"

		b1 := makeActiveRec("rt-b1", "github:44444444-4444-4444-4444-444444444444", fixed)
		b1.FamilyID = famB

		for _, r := range []*RefreshTokenRecord{a1, a2, a3, b1} {
			seedRefreshDoc(t, repo, r)
			id := r.RefreshID
			t.Cleanup(func() { _, _ = repo.docRT(id).Delete(ctx) })
		}

		n, err := repo.RevokeFamily(ctx, famA, "compromised", fixed)
		require.NoError(t, err)
		require.Equal(t, 2, n)

		gotA1 := getRefreshDoc(t, repo, "rt-a1")
		require.False(t, gotA1.RevokedAt.IsZero())
		require.Equal(t, fixed, gotA1.RevokedAt.UTC())
		require.Equal(t, "compromised", gotA1.RevokeReason)

		gotA2 := getRefreshDoc(t, repo, "rt-a2")
		require.False(t, gotA2.RevokedAt.IsZero())
		require.Equal(t, fixed, gotA2.RevokedAt.UTC())
		require.Equal(t, "compromised", gotA2.RevokeReason)

		gotA3 := getRefreshDoc(t, repo, "rt-a3")
		require.Equal(t, "already", gotA3.RevokeReason)
		require.True(t, gotA3.RevokedAt.Before(fixed))

		gotB1 := getRefreshDoc(t, repo, "rt-b1")
		require.True(t, gotB1.RevokedAt.IsZero())
		require.Equal(t, "", gotB1.RevokeReason)
	})

	t.Run("zero time uses repo.now()", func(t *testing.T) {
		t.Parallel()

		const fam = "fam-zero"
		z1 := makeActiveRec("rt-z1", "github:aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", fixed)
		z1.FamilyID = fam
		seedRefreshDoc(t, repo, z1)
		t.Cleanup(func() { _, _ = repo.docRT(z1.RefreshID).Delete(context.Background()) })

		n, err := repo.RevokeFamily(context.Background(), fam, "zero-time", time.Time{})
		require.NoError(t, err)
		require.Equal(t, 1, n)

		got := getRefreshDoc(t, repo, "rt-z1")
		require.Equal(t, fixed, got.RevokedAt.UTC())
		require.Equal(t, "zero-time", got.RevokeReason)
	})

	t.Run("invalid familyID returns error and zero count", func(t *testing.T) {
		t.Parallel()

		n, err := repo.RevokeFamily(context.Background(), "fam/invalid", "x", fixed)
		require.Error(t, err)
		require.Equal(t, 0, n)
	})
}
