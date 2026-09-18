package store

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrEmptyKMSKeyID       = errors.New("empty KMS key ID")
	ErrNilKMSStringAdapter = errors.New("nil KMS string adapter")
	ErrKMSKeyIDMismatch    = errors.New("KMS key ID mismatch")
	ErrEmptyCiphertextBlob = errors.New("empty ciphertext blob")
)

type kmsStringAdapter interface {
	EncryptString(ctx context.Context, plain string) (string, error)
	DecryptString(ctx context.Context, cipherB64 string) (string, error)
}

type KMSTokenEncryptor struct {
	kid     string
	adapter kmsStringAdapter
}

func NewKMSTokenEncryptor(
	kid string,
	adapter kmsStringAdapter,
) (*KMSTokenEncryptor, error) {
	if kid == "" {
		return nil, ErrEmptyKMSKeyID
	}
	if adapter == nil {
		return nil, ErrNilKMSStringAdapter
	}

	return &KMSTokenEncryptor{
		kid:     kid,
		adapter: adapter,
	}, nil
}

func (e *KMSTokenEncryptor) EncryptString(
	ctx context.Context,
	plain string,
) (Ciphertext, error) {
	blob, err := e.adapter.EncryptString(ctx, plain)
	if err != nil {
		return Ciphertext{}, fmt.Errorf("KMS encrypt GitHub token: %w", err)
	}

	return Ciphertext{
		KID:  e.kid,
		Blob: blob,
	}, nil
}

func (e *KMSTokenEncryptor) DecryptString(
	ctx context.Context,
	ct Ciphertext,
) (string, error) {
	if ct.KID != e.kid {
		return "", fmt.Errorf(
			"%w: got %q, want %q",
			ErrKMSKeyIDMismatch,
			ct.KID,
			e.kid,
		)
	}
	if ct.Blob == "" {
		return "", ErrEmptyCiphertextBlob
	}

	plain, err := e.adapter.DecryptString(ctx, ct.Blob)
	if err != nil {
		return "", fmt.Errorf("KMS decrypt GitHub token: %w", err)
	}

	return plain, nil
}

var _ TokenEncryptor = (*KMSTokenEncryptor)(nil)
