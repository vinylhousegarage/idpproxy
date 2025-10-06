package store

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_MarkUsed(t *testing.T) {
	requireEmulator(t)
	t.Parallel()

	fixed := time.Unix(1_800_000_000, 0)
	newRepo := func() *Repo { return newTestRepoWithNow(t, fixed) }

	type seedFn func(*testing.T, *Repo)
	mkSeed := func(rec RefreshTokenRecord) seedFn {
		return func(t *testing.T, r *Repo) { seedRefreshDoc(t, r, &rec) }
	}

	tests := []struct {
		name          string
		refreshID     string
		seed          seedFn
		wantErrIs     error
		wantErrSubstr string
		wantLastUsed  *time.Time
	}{
		{
			name:      "success",
			refreshID: "rt-ok-1",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:  "rt-ok-1",
				UserID:     "u1",
				DigestB64:  "d",
				KeyID:      "kid",
				FamilyID:   "fam1",
				CreatedAt:  fixed.Add(-time.Hour),
				ExpiresAt:  fixed.Add(time.Hour),
				DeleteAt:   time.Time{},
				RevokedAt:  time.Time{},
				LastUsedAt: time.Time{},
			}),
			wantLastUsed: &fixed,
		},
		{
			name:      "not-found",
			refreshID: "rt-does-not-exist",
			seed:      nil,
			wantErrIs: ErrNotFound,
		},
		{
			name:      "invalid-id-empty",
			refreshID: "",
			seed:      nil,
			wantErrIs: ErrInvalidID,
		},
		{
			name:      "invalid-id-has-slash",
			refreshID: "a/b",
			seed:      nil,
			wantErrIs: ErrInvalidID,
		},
		{
			name:      "revoked",
			refreshID: "rt-revoked-1",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:  "rt-revoked-1",
				UserID:     "u1",
				DigestB64:  "d",
				KeyID:      "kid",
				FamilyID:   "fam1",
				CreatedAt:  fixed.Add(-2 * time.Hour),
				ExpiresAt:  fixed.Add(time.Hour),
				RevokedAt:  fixed.Add(-30 * time.Minute),
				LastUsedAt: time.Time{},
			}),
			wantErrSubstr: "token revoked",
		},
		{
			name:      "expired",
			refreshID: "rt-expired-1",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:  "rt-expired-1",
				UserID:     "u1",
				DigestB64:  "d",
				KeyID:      "kid",
				FamilyID:   "fam1",
				CreatedAt:  fixed.Add(-48 * time.Hour),
				ExpiresAt:  fixed.Add(-1 * time.Minute),
				LastUsedAt: time.Time{},
			}),
			wantErrSubstr: "token expired",
		},
		{
			name:      "deleted",
			refreshID: "rt-deleted-1",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:  "rt-deleted-1",
				UserID:     "u1",
				DigestB64:  "d",
				KeyID:      "kid",
				FamilyID:   "fam1",
				CreatedAt:  fixed.Add(-48 * time.Hour),
				ExpiresAt:  fixed.Add(time.Hour),
				DeleteAt:   fixed.Add(-1 * time.Minute),
				LastUsedAt: time.Time{},
			}),
			wantErrSubstr: "token deleted",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRepo()
			if tt.seed != nil {
				tt.seed(t, r)
			}

			err := r.MarkUsed(context.Background(), tt.refreshID)
			if tt.wantErrIs != nil {
				require.ErrorIs(t, err, tt.wantErrIs)
				return
			}
			if tt.wantErrSubstr != "" {
				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), tt.wantErrSubstr), "got err: %v", err)
				return
			}

			require.NoError(t, err)
			if tt.wantLastUsed != nil {
				got := getRefreshDoc(t, r, tt.refreshID)
				require.True(t, got.LastUsedAt.Equal(*tt.wantLastUsed),
					"expected %v, got %v", tt.wantLastUsed, got.LastUsedAt)
			}
		})
	}
}
