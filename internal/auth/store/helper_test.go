package store

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func newTestRepoWithNow(t *testing.T, fixed time.Time) *Repo {
	t.Helper()
	r := newTestRepo(t)
	r.now = func() time.Time { return fixed }
	return r
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

func seedRefreshDoc(t *testing.T, r *Repo, rec *RefreshTokenRecord) {
	t.Helper()
	ctx := context.Background()

	_, err := r.docRT(rec.RefreshID).Create(ctx, rec)
	if err == nil {
		return
	}

	if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
		_, err = r.docRT(rec.RefreshID).Set(ctx, rec)
	}

	require.NoError(t, err)
}

func getRefreshDoc(t *testing.T, r *Repo, id string) *RefreshTokenRecord {
	t.Helper()
	ctx := context.Background()
	snap, err := r.docRT(id).Get(ctx)
	require.NoError(t, err)

	var rec RefreshTokenRecord
	require.NoError(t, snap.DataTo(&rec))
	return &rec
}

func makeActiveRec(id, user string, now time.Time) *RefreshTokenRecord {
	rec := makeRec(id, user, "fam-1", now)
	rec.CreatedAt = now.Add(-time.Hour)
	rec.ExpiresAt = now.Add(24 * time.Hour)
	rec.DeleteAt = now.Add(30 * 24 * time.Hour)
	rec.RevokedAt = time.Time{}
	rec.ReplacedBy = ""

	return rec
}
