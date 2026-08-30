package store

import (
	"context"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFirestoreGitHubTokenRepo_DeleteByFirebaseUID(t *testing.T) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("FIRESTORE_EMULATOR_HOST is not set")
	}

	t.Parallel()

	ctx := context.Background()
	client := newFirestoreTestClient(t, ctx)

	const firebaseUID = "delete-test-firebase-uid"

	doc := client.
		Collection(collectionGitHubTokens).
		Doc(firebaseUID)

	_, err := doc.Set(ctx, firestoreGitHubToken{
		FirebaseUID: firebaseUID,
		GitHubID:    "github-user-id",
		Provider:    "github",
		Login:       "test-user",
	})
	if err != nil {
		t.Fatalf("create GitHub token: %v", err)
	}

	repo := NewFirestoreGitHubTokenRepo(
		client,
		fakeTokenEncryptor{},
	)

	if err := repo.DeleteByFirebaseUID(ctx, firebaseUID); err != nil {
		t.Fatalf("delete GitHub token: %v", err)
	}

	_, err = doc.Get(ctx)
	if status.Code(err) != codes.NotFound {
		t.Fatalf(
			"expected deleted document to return NotFound, got: %v",
			err,
		)
	}

	if err := repo.DeleteByFirebaseUID(ctx, firebaseUID); err != nil {
		t.Fatalf("delete missing GitHub token: %v", err)
	}
}
