package kms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"cloud.google.com/go/kms/apiv1/kmspb"
)

type Adapter struct {
	c      KMSClient
	keyRes string
	aad    []byte
}

func NewAdapter(c KMSClient, keyRes string, aad []byte) (*Adapter, error) {
	if c == nil {
		return nil, errors.New("kms: nil client")
	}
	if keyRes == "" {
		return nil, errors.New("kms: empty key resource")
	}

	return &Adapter{c: c, keyRes: keyRes, aad: aad}, nil
}

func (a *Adapter) EncryptString(ctx context.Context, plain string) (string, error) {
	req := &kmspb.EncryptRequest{
		Name:                        a.keyRes,
		Plaintext:                   []byte(plain),
		AdditionalAuthenticatedData: a.aad,
	}

	resp, err := a.c.Encrypt(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrEncryptFailed, err)
	}

	return base64.StdEncoding.EncodeToString(resp.Ciphertext), nil
}

func (a *Adapter) DecryptString(ctx context.Context, cipherB64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrBadFormat, err)
	}

	req := &kmspb.DecryptRequest{
		Name:                        a.keyRes,
		Ciphertext:                  ciphertext,
		AdditionalAuthenticatedData: a.aad,
	}

	resp, err := a.c.Decrypt(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDecryptFailed, err)
	}

	return string(resp.Plaintext), nil
}
