package store

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestFirestoreGitHubTokenRepo_Upsert(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	now := time.Date(2026, time.August, 23, 0, 0, 0, 0, time.UTC)

	encryptedToken := Ciphertext{
		KID:  "github-token-key-v1",
		Blob: "encrypted-access-token",
	}

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{
			ciphertext: encryptedToken,
		},
	)
	repo.now = func() time.Time {
		return now
	}

	firebaseUID := fmt.Sprintf("firebase-user-%d", time.Now().UnixNano())

	t.Cleanup(func() {
		_, _ = repo.col.Doc(firebaseUID).Delete(context.Background())
	})

	rec := &GitHubTokenRecord{
		GitHubID:    "123456",
		Provider:    "github",
		FirebaseUID: firebaseUID,
		Login:       "vinylhousegarage",
		Scopes:      []string{"read:user", "repo"},
		TokenType:   "bearer",
		AccessToken: "plain-access-token",
		ExpiresAt:   now.Add(24 * time.Hour),
		LastUsedAt:  now.Add(-1 * time.Hour),
	}

	if err := repo.Upsert(ctx, rec); err != nil {
		t.Fatalf("upsert GitHub token: %v", err)
	}

	snapshot, err := repo.col.Doc(rec.FirebaseUID).Get(ctx)
	if err != nil {
		t.Fatalf("get saved GitHub token: %v", err)
	}

	var got firestoreGitHubToken
	if err := snapshot.DataTo(&got); err != nil {
		t.Fatalf("decode saved GitHub token: %v", err)
	}

	want := firestoreGitHubToken{
		GitHubID:    rec.GitHubID,
		Provider:    rec.Provider,
		FirebaseUID: rec.FirebaseUID,
		Login:       rec.Login,
		Scopes:      rec.Scopes,
		TokenType:   rec.TokenType,
		AccessToken: encryptedToken,
		ExpiresAt:   rec.ExpiresAt,
		LastUsedAt:  rec.LastUsedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"saved token mismatch (-got +want):\n got: %#v\nwant: %#v",
			got,
			want,
		)
	}
}
