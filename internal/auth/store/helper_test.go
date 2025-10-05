package store

import (
	"context"
	"os"
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

func makeRec(id, user, family string, now time.Time) *RefreshTokenRecord {
	return &RefreshTokenRecord{
		RefreshID:    id,
		UserID:       user,
		DigestB64:    "dGVzdC1kaWdlc3Q=",
		KeyID:        "kid-1",
		FamilyID:     family,
		ReplacedBy:   "",
		RevokedAt:    time.Time{},
		RevokeReason: "",
		CreatedAt:    now,
		LastUsedAt:   time.Time{},
		ExpiresAt:    now.Add(24 * time.Hour),
		DeleteAt:     now.Add(48 * time.Hour),
	}
}
