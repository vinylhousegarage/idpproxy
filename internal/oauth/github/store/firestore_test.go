package store

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
)

type fakeTokenEncryptor struct {
	ciphertext Ciphertext
	plaintext  string
	encryptErr error
	decryptErr error
}

func (e fakeTokenEncryptor) EncryptString(_ string) (Ciphertext, error) {
	if e.encryptErr != nil {
		return Ciphertext{}, e.encryptErr
	}

	return e.ciphertext, nil
}

func (e fakeTokenEncryptor) DecryptString(_ Ciphertext) (string, error) {
	if e.decryptErr != nil {
		return "", e.decryptErr
	}

	return e.plaintext, nil
}

func newFirestoreTestClient(
	t *testing.T,
	ctx context.Context,
) *firestore.Client {
	t.Helper()

	client, err := firestore.NewClient(ctx, "idpproxy-test")
	if err != nil {
		t.Fatalf("create Firestore client: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}
