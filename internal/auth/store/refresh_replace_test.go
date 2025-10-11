package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

	mkID := func(t *testing.T, suffix string) string {
		base := strings.ReplaceAll(t.Name(), "/", "_")
		return fmt.Sprintf("%s-%d-%s", base, time.Now().UnixNano(), suffix)
	}

	cleanupIDs := func(t *testing.T, repo *Repo, ids ...string) {
		t.Helper()
		t.Cleanup(func() {
			for _, id := range ids {
				if id == "" {
					continue
				}
				_, _ = repo.docRT(id).Delete(context.Background())
			}
		})
	}

	t.Run("success: old active → old closed & new created", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		oldID := mkID(t, "old")
		newID := mkID(t, "new")
		cleanupIDs(t, repo, oldID, newID)

		old := makeActiveRec(oldID, userOK, now)
		seedRefreshDoc(t, repo, old)

		newRec := &RefreshTokenRecord{
			RefreshID: newID,
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(48 * time.Hour),
			DeleteAt:  now.Add(60 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, old.RefreshID, newRec, now)
		require.NoError(t, err)

		oldSnap, err := repo.docRT(oldID).Get(ctx)
		require.NoError(t, err)
		var gotOld RefreshTokenRecord
		require.NoError(t, oldSnap.DataTo(&gotOld))
		require.Equal(t, newID, gotOld.ReplacedBy)
		require.Equal(t, now, gotOld.RevokedAt)

		newSnap, err := repo.docRT(newID).Get(ctx)
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

		newID := mkID(t, "new404")
		cleanupIDs(t, repo, newID)

		newRec := &RefreshTokenRecord{
			RefreshID: newID,
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, mkID(t, "missing-old"), newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNotFound), "got %v", err)
	})

	t.Run("invalid args: nil newRec → ErrInvalid", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		err := repo.Replace(ctx, mkID(t, "any"), nil, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalid))
	})

	t.Run("invalid args: bad oldID (contains slash) → ErrInvalidID (wrapped)", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		newID := mkID(t, "new-bad-old")
		cleanupIDs(t, repo, newID)

		newRec := &RefreshTokenRecord{
			RefreshID: newID,
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

		oldID := mkID(t, "old-u")
		newID := mkID(t, "new-u")
		cleanupIDs(t, repo, oldID, newID)

		old := makeActiveRec(oldID, userOK, now)
		seedRefreshDoc(t, repo, old)

		newRec := &RefreshTokenRecord{
			RefreshID: newID,
			UserID:    userNG,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, oldID, newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrConflict))
	})

	t.Run("conflict: old is already revoked (inactive) → ErrConflict", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		oldID := mkID(t, "old-r")
		newID := mkID(t, "new-r")
		cleanupIDs(t, repo, oldID, newID)

		old := makeActiveRec(oldID, userOK, now)
		old.RevokedAt = now.Add(-time.Minute)
		seedRefreshDoc(t, repo, old)

		newRec := &RefreshTokenRecord{
			RefreshID: newID,
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, oldID, newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrConflict))
	})

	t.Run("conflict: new RefreshID already exists → ErrConflict", func(t *testing.T) {
		t.Parallel()
		repo := newRepo()

		oldID := mkID(t, "old-dup")
		newID := mkID(t, "new-dup")
		cleanupIDs(t, repo, oldID, newID)

		old := makeActiveRec(oldID, userOK, now)
		seedRefreshDoc(t, repo, old)

		existingNew := makeActiveRec(newID, userOK, now)
		seedRefreshDoc(t, repo, existingNew)

		newRec := &RefreshTokenRecord{
			RefreshID: newID,
			UserID:    userOK,
			DigestB64: "ZGVtbw==",
			KeyID:     "k1",
			ExpiresAt: now.Add(24 * time.Hour),
			DeleteAt:  now.Add(30 * 24 * time.Hour),
		}

		err := repo.Replace(ctx, oldID, newRec, now)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrConflict), "got %v", err)
	})
}
