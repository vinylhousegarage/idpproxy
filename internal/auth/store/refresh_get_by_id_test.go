package store

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

func requireEmulator(t *testing.T) {
	t.Helper()
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set; skipping Firestore emulator tests")
	}
}

func newTestRepo(t *testing.T) *Repo {
	t.Helper()

	ctx := context.Background()
	projectID := os.Getenv("TEST_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Fatal("TEST_FIRESTORE_PROJECT is not set")
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithoutAuthentication())
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	fixed := time.Unix(1_725_000_000, 0)
	return &Repo{fs: client, now: func() time.Time { return fixed }}
}

func makeRec(id, user, fam string, now time.Time) *RefreshTokenRecord {
	return &RefreshTokenRecord{
		RefreshID:    id,
		UserID:       user,
		DigestB64:    "dGVzdC1kaWdlc3Q=",
		KeyID:        "kid-1",
		FamilyID:     fam,
		ReplacedBy:   "",
		RevokedAt:    time.Time{},
		CreatedAt:    now,
		LastUsedAt:   time.Time{},
		ExpiresAt:    now.Add(24 * time.Hour),
		DeleteAt:     now.Add(48 * time.Hour),
		RevokeReason: "",
	}
}

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
		errSubstr  string
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
			name:      "invalid-empty",
			id:        "",
			wantErr:   true,
			errSubstr: "empty refreshid",
		},
		{
			name:      "invalid-slash",
			id:        "bad/id",
			wantErr:   true,
			errSubstr: "must not contain '/'",
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
				}
				if tt.errSubstr != "" {
					require.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errSubstr))
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
