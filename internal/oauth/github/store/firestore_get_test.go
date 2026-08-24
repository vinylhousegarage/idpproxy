package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestFirestoreGitHubTokenRepo_GetByFirebaseUID(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	firebaseUID := fmt.Sprintf(
		"firebase-uid-get-test-%d",
		time.Now().UnixNano(),
	)

	now := time.Date(2026, time.August, 24, 0, 0, 0, 0, time.UTC)

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{
			plaintext: "decrypted-access-token",
		},
	)

	doc := repo.col.Doc(firebaseUID)
	t.Cleanup(func() {
		_, _ = doc.Delete(context.Background())
	})

	_, err := doc.Set(ctx, firestoreGitHubToken{
		GitHubID:    "123456",
		Provider:    "github",
		FirebaseUID: firebaseUID,
		Login:       "vinylhousegarage",
		Scopes:      []string{"read:user", "repo"},
		TokenType:   "bearer",
		AccessToken: Ciphertext{},
		ExpiresAt:   now.Add(time.Hour),
		LastUsedAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("set GitHub token document: %v", err)
	}

	got, err := repo.GetByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		t.Fatalf("GetByFirebaseUID() error = %v", err)
	}

	want := &GitHubTokenRecord{
		GitHubID:    "123456",
		Provider:    "github",
		FirebaseUID: firebaseUID,
		Login:       "vinylhousegarage",
		Scopes:      []string{"read:user", "repo"},
		TokenType:   "bearer",
		AccessToken: "decrypted-access-token",
		ExpiresAt:   now.Add(time.Hour),
		LastUsedAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetByFirebaseUID() = %+v, want %+v", got, want)
	}
}

func TestFirestoreGitHubTokenRepo_GetByFirebaseUID_NotFound(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{},
	)

	firebaseUID := fmt.Sprintf(
		"firebase-uid-get-not-found-test-%d",
		time.Now().UnixNano(),
	)

	got, err := repo.GetByFirebaseUID(ctx, firebaseUID)
	if !errors.Is(err, ErrGitHubTokenNotFound) {
		t.Errorf(
			"GetByFirebaseUID() error = %v, want ErrGitHubTokenNotFound",
			err,
		)
	}

	if got != nil {
		t.Errorf("GetByFirebaseUID() = %+v, want nil", got)
	}
}

func TestFirestoreGitHubTokenRepo_GetByFirebaseUID_DecryptError(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	firebaseUID := fmt.Sprintf(
		"firebase-uid-get-decrypt-error-test-%d",
		time.Now().UnixNano(),
	)

	decryptErr := errors.New("decrypt failed")

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{
			decryptErr: decryptErr,
		},
	)

	doc := repo.col.Doc(firebaseUID)
	t.Cleanup(func() {
		_, _ = doc.Delete(context.Background())
	})

	_, err := doc.Set(ctx, firestoreGitHubToken{
		FirebaseUID: firebaseUID,
		AccessToken: Ciphertext{},
	})
	if err != nil {
		t.Fatalf("set GitHub token document: %v", err)
	}

	got, err := repo.GetByFirebaseUID(ctx, firebaseUID)
	if !errors.Is(err, decryptErr) {
		t.Errorf(
			"GetByFirebaseUID() error = %v, want wrapped %v",
			err,
			decryptErr,
		)
	}

	if got != nil {
		t.Errorf("GetByFirebaseUID() = %+v, want nil", got)
	}
}
