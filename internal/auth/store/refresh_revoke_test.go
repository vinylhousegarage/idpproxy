package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepo_Revoke(t *testing.T) {
	requireEmulator(t)
	t.Parallel()

	fixed := time.Unix(1_800_000_000, 0)
	newRepo := func() *Repo { return newTestRepoWithNow(t, fixed) }
	ctx := context.Background()

	type seedFn func(*testing.T, *Repo)
	mkSeed := func(rec RefreshTokenRecord) seedFn {
		return func(t *testing.T, r *Repo) { seedRefreshDoc(t, r, &rec) }
	}

	tests := []struct {
		name          string
		refreshID     string
		seed          seedFn
		reason        string
		wantErrIs     error
		wantErrSubstr string
		wantRevokedAt *time.Time
	}{
		{
			name:      "success-active-token",
			refreshID: "rt-ok-1",
			reason:    "manual-logout",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:  "rt-ok-1",
				UserID:     "u1",
				DigestB64:  "d1",
				KeyID:      "k1",
				FamilyID:   "fam-1",
				CreatedAt:  fixed.Add(-1 * time.Hour),
				ExpiresAt:  fixed.Add(24 * time.Hour),
				DeleteAt:   fixed.Add(48 * time.Hour),
				LastUsedAt: time.Time{},
			}),
			wantRevokedAt: &fixed,
		},
		{
			name:      "already-revoked",
			refreshID: "rt-revoked-1",
			reason:    "dup",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:    "rt-revoked-1",
				UserID:       "u1",
				DigestB64:    "d2",
				KeyID:        "k1",
				FamilyID:     "fam-1",
				CreatedAt:    fixed.Add(-2 * time.Hour),
				ExpiresAt:    fixed.Add(24 * time.Hour),
				DeleteAt:     fixed.Add(48 * time.Hour),
				RevokedAt:    fixed.Add(-1 * time.Minute),
				RevokeReason: "prior",
			}),
			wantErrIs: ErrAlreadyRevoked,
		},
		{
			name:      "already-replaced",
			refreshID: "rt-old-1",
			reason:    "rotate",
			seed: mkSeed(RefreshTokenRecord{
				RefreshID:  "rt-old-1",
				UserID:     "u1",
				DigestB64:  "d3",
				KeyID:      "k1",
				FamilyID:   "fam-2",
				CreatedAt:  fixed.Add(-2 * time.Hour),
				ExpiresAt:  fixed.Add(24 * time.Hour),
				DeleteAt:   fixed.Add(48 * time.Hour),
				ReplacedBy: "rt-new-1",
			}),
			wantErrIs: ErrAlreadyRevoked,
		},
		{
			name:          "not-found",
			refreshID:     "rt-no-such-id",
			reason:        "any",
			seed:          nil,
			wantErrIs:     ErrNotFound,
			wantErrSubstr: "not found",
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

			err := r.Revoke(ctx, tt.refreshID, tt.reason, fixed)
			if tt.wantErrIs != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErrIs)
				if tt.wantErrSubstr != "" {
					require.Contains(t, err.Error(), tt.wantErrSubstr)
				}
				return
			}
			require.NoError(t, err)

			got, err := r.GetByID(ctx, tt.refreshID)
			require.NoError(t, err)
			require.Equal(t, tt.reason, got.RevokeReason)
			if tt.wantRevokedAt != nil {
				require.WithinDuration(t, *tt.wantRevokedAt, got.RevokedAt, time.Second)
			}
		})
	}
}
