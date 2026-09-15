package store

import "context"

type Ciphertext struct {
	KID  string `firestore:"kid"`
	Blob string `firestore:"blob"`
}

type TokenEncryptor interface {
	EncryptString(
		ctx context.Context,
		plain string,
	) (Ciphertext, error)

	DecryptString(
		ctx context.Context,
		ct Ciphertext,
	) (string, error)
}
