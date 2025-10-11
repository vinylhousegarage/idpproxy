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

	fixed := time.Unix(1_800_000_000, 0).UTC()
	newRepo := func() *Repo { return newTestRepoWithNow(t, fixed) }

	mkID := func(t *testing.T, suffix string) string {
		base := strings.ReplaceAll(t.Name(), "/", "_")
		return base + "-" + time.Now().UTC().Format("20060102T150405.000000000Z") + "-" + suffix
	}

	cleanupIDs := func(t *testing.T, r *Repo, ids ...string) {
		t.Helper()
		t.Cleanup(func() {
			for _, id := range ids {
				if id == "" {
					continue
				}
				_, _ = r.docRT(id).Delete(context.Background())
			}
		})
	}

	type seedFn func(*testing.T, *Repo, string)

	seedActive := func(user string) seedFn {
		return func(t *testing.T, r *Repo, id string) {
			rec := makeActiveRec(id, user, fixed)
			seedRefreshDoc(t, r, rec)
		}
	}
	seedWith := func(mut func(*RefreshTokenRecord)) seedFn {
		return func(t *testing.T, r *Repo, id string) {
			rec := makeActiveRec(id, "github:123e4567-e89b-12d3-a456-426614174000", fixed)
			mut(rec)
			seedRefreshDoc(t, r, rec)
		}
	}

	tests := []struct {
		name          string
		seed          seedFn
		refreshID     string
		wantErrIs     error
		wantErrSubstr string
		wantLastUsed  *time.Time
	}{
		{
			name:         "success",
			seed:         seedActive("github:123e4567-e89b-12d3-a456-426614174000"),
			wantLastUsed: &fixed,
		},
		{
			name:      "not-found",
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
			name: "revoked",
			seed: seedWith(func(r *RefreshTokenRecord) {
				r.RevokedAt = fixed.Add(-30 * time.Minute)
			}),
			wantErrSubstr: "token revoked",
		},
		{
			name: "expired",
			seed: seedWith(func(r *RefreshTokenRecord) {
				r.ExpiresAt = fixed.Add(-1 * time.Minute)
			}),
			wantErrSubstr: "token expired",
		},
		{
			name: "deleted",
			seed: seedWith(func(r *RefreshTokenRecord) {
				r.DeleteAt = fixed.Add(-1 * time.Minute)
			}),
			wantErrSubstr: "token deleted",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRepo()

			var id string
			if tt.refreshID != "" {
				id = tt.refreshID
			} else {
				id = mkID(t, "rt")
			}

			if tt.seed != nil && tt.refreshID == "" {
				cleanupIDs(t, r, id)
				tt.seed(t, r, id)
			}

			err := r.MarkUsed(context.Background(), id)

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
				got := getRefreshDoc(t, r, id)
				require.True(t, got.LastUsedAt.Equal(*tt.wantLastUsed),
					"expected %v, got %v", *tt.wantLastUsed, got.LastUsedAt)
			}
		})
	}
}
