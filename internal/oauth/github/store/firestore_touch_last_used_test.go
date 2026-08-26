package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFirestoreGitHubTokenRepo_TouchLastUsed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	before := time.Date(2026, 8, 27, 0, 0, 0, 0, time.UTC)
	now := before.Add(time.Hour)

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{},
	)
	repo.now = func() time.Time {
		return now
	}

	firebaseUID := "touch-last-used-test"

	_, err := client.Collection(collectionGitHubTokens).Doc(firebaseUID).Set(
		ctx,
		firestoreGitHubToken{
			FirebaseUID: firebaseUID,
			LastUsedAt:  before,
			UpdatedAt:   before,
		},
	)
	if err != nil {
		t.Fatalf("set GitHub token: %v", err)
	}

	if err := repo.TouchLastUsed(ctx, firebaseUID); err != nil {
		t.Fatalf("touch last used: %v", err)
	}

	snapshot, err := client.Collection(collectionGitHubTokens).Doc(firebaseUID).Get(ctx)
	if err != nil {
		t.Fatalf("get GitHub token: %v", err)
	}

	var got firestoreGitHubToken
	if err := snapshot.DataTo(&got); err != nil {
		t.Fatalf("decode GitHub token: %v", err)
	}

	if !got.LastUsedAt.Equal(now) {
		t.Errorf("LastUsedAt = %v, want %v", got.LastUsedAt, now)
	}

	if !got.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
	}
}

func TestFirestoreGitHubTokenRepo_TouchLastUsed_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{},
	)

	err := repo.TouchLastUsed(ctx, "not-found-touch-last-used-test")
	if !errors.Is(err, ErrGitHubTokenNotFound) {
		t.Errorf("TouchLastUsed() error = %v, want %v", err, ErrGitHubTokenNotFound)
	}
}
