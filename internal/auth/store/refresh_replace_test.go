package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_Replace(t *testing.T) {
	requireEmulator(t)

	t.Parallel()
	ctx := context.Background()
	now := time.Unix(1_800_000_000, 0).UTC()

	newRepo := func() *Repo { return newTestRepo(t) }

	const userOK = "github:123e4567-e89b-12d3-a456-426614174000"
	const userNG = "github:00000000-0000-0000-0000-000000000000"

	t.Run("success: old active → old closed & new created", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		old := makeActiveRec("rt-old-1", userOK, now)
		seedRefreshDoc(t, repo, old)

		newRec := &RefreshTokenRecord{
			RefreshID: "rt-new-1",
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(48 * time.Hour),
			DeleteAt:  now.Add(60 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, old.RefreshID, newRec, now)
		require.NoError(t, err)

		oldSnap, err := repo.docRT(old.RefreshID).Get(ctx)
		require.NoError(t, err)

		var gotOld RefreshTokenRecord
		require.NoError(t, oldSnap.DataTo(&gotOld))
		require.Equal(t, "rt-new-1", gotOld.ReplacedBy)
		require.Equal(t, now, gotOld.RevokedAt)

		newSnap, err := repo.docRT(newRec.RefreshID).Get(ctx)
		require.NoError(t, err)

		var gotNew RefreshTokenRecord
		require.NoError(t, newSnap.DataTo(&gotNew))
		require.Equal(t, old.FamilyID, gotNew.FamilyID)
		require.Equal(t, now, gotNew.CreatedAt)
		require.True(t, gotNew.RevokedAt.IsZero())
		require.Equal(t, "", gotNew.ReplacedBy)
	})

	t.Run("not found: oldID does not exist → ErrNotFound", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		newRec := &RefreshTokenRecord{
			RefreshID: "rt-new-404",
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, "rt-missing", newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNotFound), "got %v", err)
	})

	t.Run("invalid args: nil newRec → ErrInvalid", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		err := repo.Replace(ctx, "rt-any", nil, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalid))
	})

	t.Run("invalid args: bad oldID (contains slash) → ErrInvalidID (wrapped)", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		newRec := &RefreshTokenRecord{
			RefreshID: "rt-new-bad-old",
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, "bad/old", newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidID), "want ErrInvalidID wrapped, got %v", err)
	})

	t.Run("conflict: user mismatch → ErrConflict", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		old := makeActiveRec("rt-old-u", userOK, now)
		seedRefreshDoc(t, repo, old)

		newRec := &RefreshTokenRecord{
			RefreshID: "rt-new-u",
			UserID:    userNG,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, old.RefreshID, newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrConflict))
	})

	t.Run("conflict: old is already revoked (inactive) → ErrConflict", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		old := makeActiveRec("rt-old-r", userOK, now)
		old.RevokedAt = now.Add(-time.Minute)
		seedRefreshDoc(t, repo, old)

		newRec := &RefreshTokenRecord{
			RefreshID: "rt-new-r",
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, old.RefreshID, newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrConflict))
	})

	t.Run("conflict: new RefreshID already exists → ErrConflict", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		old := makeActiveRec("rt-old-dup", userOK, now)
		seedRefreshDoc(t, repo, old)

		existingNew := makeActiveRec("rt-new-dup", userOK, now)
		seedRefreshDoc(t, repo, existingNew)

		newRec := &RefreshTokenRecord{
			RefreshID: "rt-new-dup",
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, old.RefreshID, newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrConflict))
	})
}
