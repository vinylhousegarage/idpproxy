package store

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepo_GetByID(t *testing.T) {
	t.Parallel()
	requireEmulator(t)

	ctx := context.Background()
	r := newTestRepo(t)
	now := r.now()

	tests := []struct {
		name       string
		seed       *RefreshTokenRecord
		id         string
		wantErr    bool
		isNotFound bool
	}{
		{
			name: "success",
			seed: makeRec("rt-1", "user-1", "fam-1", now),
			id:   "rt-1",
		},
		{
			name:       "not-found",
			id:         "does-not-exist",
			wantErr:    true,
			isNotFound: true,
		},
		{
			name:    "invalid-empty",
			id:      "",
			wantErr: true,
		},
		{
			name:    "invalid-slash",
			id:      "bad/id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.seed != nil {
				tt.seed.RefreshID = tt.seed.RefreshID + "-" + strings.ReplaceAll(t.Name(), "/", "_")
				tt.id = tt.seed.RefreshID

				_, err := r.docRT(tt.seed.RefreshID).Set(ctx, tt.seed)
				require.NoError(t, err)
				t.Cleanup(func() { _, _ = r.docRT(tt.seed.RefreshID).Delete(ctx) })
			}

			got, err := r.GetByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				if tt.isNotFound {
					require.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound, got %v", err)
				} else {
					if ErrInvalidID != nil {
						require.True(t, errors.Is(err, ErrInvalidID), "expected ErrInvalidID, got: %v", err)
					}
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			require.Equal(t, tt.seed.RefreshID, got.RefreshID)
			require.Equal(t, tt.seed.UserID, got.UserID)
			require.Equal(t, tt.seed.DigestB64, got.DigestB64)
			require.Equal(t, tt.seed.KeyID, got.KeyID)
			require.Equal(t, tt.seed.FamilyID, got.FamilyID)
			require.True(t, got.CreatedAt.Equal(tt.seed.CreatedAt))
			require.True(t, got.ExpiresAt.Equal(tt.seed.ExpiresAt))
			require.True(t, got.DeleteAt.Equal(tt.seed.DeleteAt))
		})
	}
}
