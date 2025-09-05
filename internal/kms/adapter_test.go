package kms

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/require"
)

type fakeKMS struct {
	encryptErr error
	decryptErr error

	lastEncryptReq *kmspb.EncryptRequest
	lastDecryptReq *kmspb.DecryptRequest
}

func (f *fakeKMS) Encrypt(ctx context.Context, req *kmspb.EncryptRequest, _ ...gax.CallOption) (*kmspb.EncryptResponse, error) {
    f.lastEncryptReq = req
    if f.encryptErr != nil {
        return nil, f.encryptErr
    }
    c := append([]byte("CIPH|"), req.Plaintext...)
    return &kmspb.EncryptResponse{Ciphertext: c}, nil
}

func (f *fakeKMS) Decrypt(ctx context.Context, req *kmspb.DecryptRequest, _ ...gax.CallOption) (*kmspb.DecryptResponse, error) {
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

		adp, err := NewAdapter(kmsFake, key, aad)
		require.NoError(t, err)

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
		adp, err := NewAdapter(kmsFake, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)
		require.NoError(t, err)

		_, err = adp.DecryptString(context.Background(), "not-base64!!!!")
		require.ErrorIs(t, err, ErrBadFormat)
	})

	t.Run("Encrypt_KMSError", func(t *testing.T) {
		t.Parallel()

		kmsFake := &fakeKMS{encryptErr: errors.New("kms down")}
		adp, err := NewAdapter(kmsFake, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)
		require.NoError(t, err)

		_, err = adp.EncryptString(context.Background(), "data")
		require.ErrorIs(t, err, ErrEncryptFailed)
	})

	t.Run("Decrypt_KMSError", func(t *testing.T) {
		t.Parallel()

		ok := &fakeKMS{}
		adpOK, err := NewAdapter(ok, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)
		require.NoError(t, err)

		c, err := adpOK.EncryptString(context.Background(), "data")
		require.NoError(t, err)

		bad := &fakeKMS{decryptErr: errors.New("kms down")}
		adpBad, err := NewAdapter(bad, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)
		require.NoError(t, err)

		_, err = adpBad.DecryptString(context.Background(), c)
		require.ErrorIs(t, err, ErrDecryptFailed)
	})

	t.Run("EncryptDecrypt_EmptyString", func(t *testing.T) {
		t.Parallel()

		kmsFake := &fakeKMS{}
		adp, err := NewAdapter(kmsFake, "projects/p/locations/l/keyRings/r/cryptoKeys/k", nil)
		require.NoError(t, err)

		enc, err := adp.EncryptString(context.Background(), "")
		require.NoError(t, err)
		require.NotEmpty(t, enc)

		out, err := adp.DecryptString(context.Background(), enc)
		require.NoError(t, err)
		require.Equal(t, "", out)
	})
}
