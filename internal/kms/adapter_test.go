package kms

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/stretchr/testify/require"
)

type fakeKMS struct {
	encryptErr error
	decryptErr error

	lastEncryptReq *kmspb.EncryptRequest
	lastDecryptReq *kmspb.DecryptRequest
}

func (f *fakeKMS) Encrypt(ctx context.Context, req *kmspb.EncryptRequest, _ ...any) (*kmspb.EncryptResponse, error) {
	f.lastEncryptReq = req
	if f.encryptErr != nil {
		return nil, f.encryptErr
	}
	c := append([]byte("CIPH|"), req.Plaintext...)
	return &kmspb.EncryptResponse{Ciphertext: c}, nil
}

func (f *fakeKMS) Decrypt(ctx context.Context, req *kmspb.DecryptRequest, _ ...any) (*kmspb.DecryptResponse, error) {
	f.lastDecryptReq = req
	if f.decryptErr != nil {
		return nil, f.decryptErr
	}
	const p = "CIPH|"
	if !strings.HasPrefix(string(req.Ciphertext), p) {
		return nil, errors.New("ciphertext malformed")
	}
	plain := []byte(string(req.Ciphertext)[len(p):])
	return &kmspb.DecryptResponse{Plaintext: plain}, nil
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run("EncryptDecrypt_Success", func(t *testing.T) {
		t.Parallel()

		kmsFake := &fakeKMS{}
		key := "projects/p/locations/global/keyRings/r/cryptoKeys/k"
		aad := []byte("uid:12345")
		adp := NewAdapter(kmsFake, key, aad)

		plain := "hello-world"
		ctx := context.Background()

		encB64, err := adp.EncryptString(ctx, plain)
		require.NoError(t, err)
		require.NotEmpty(t, encB64)

		require.Equal(t, key, kmsFake.lastEncryptReq.GetName())
		require.Equal(t, aad, kmsFake.lastEncryptReq.GetAdditionalAuthenticatedData())

		cipherBytes, err := base64.StdEncoding.DecodeString(encB64)
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(string(cipherBytes), "CIPH|"))

		out, err := adp.DecryptString(ctx, encB64)
		require.NoError(t, err)
		require.Equal(t, plain, out)

		require.Equal(t, key, kmsFake.lastDecryptReq.GetName())
		require.Equal(t, aad, kmsFake.lastDecryptReq.GetAdditionalAuthenticatedData())
	})

	t.Run("DecryptString_BadBase64", func(t *testing.T) {
		t.Parallel()

		kmsFake := &fakeKMS{}
		adp := NewAdapter(kmsFake, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)

		_, err := adp.DecryptString(context.Background(), "not-base64!!!!")
		require.Error(t, err)
		require.Contains(t, err.Error(), "bad ciphertext format")
	})

	t.Run("Encrypt_KMSError", func(t *testing.T) {
		t.Parallel()

		kmsFake := &fakeKMS{encryptErr: errors.New("kms down")}
		adp := NewAdapter(kmsFake, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)

		_, err := adp.EncryptString(context.Background(), "data")
		require.Error(t, err)
		require.Contains(t, err.Error(), "kms encrypt")
	})

	t.Run("Decrypt_KMSError", func(t *testing.T) {
		t.Parallel()

		ok := &fakeKMS{}
		adpOK := NewAdapter(ok, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)
		c, err := adpOK.EncryptString(context.Background(), "data")
		require.NoError(t, err)

		bad := &fakeKMS{decryptErr: errors.New("kms down")}
		adpBad := NewAdapter(bad, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)

		_, err = adpBad.DecryptString(context.Background(), c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "kms decrypt")
	})

	t.Run("EncryptDecrypt_EmptyString", func(t *testing.T) {
		t.Parallel()

		kmsFake := &fakeKMS{}
		adp := NewAdapter(kmsFake, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)

		enc, err := adp.EncryptString(context.Background(), "")
		require.NoError(t, err)
		require.NotEmpty(t, enc)

		out, err := adp.DecryptString(context.Background(), enc)
		require.NoError(t, err)
		require.Equal(t, "", out)
	})
}
