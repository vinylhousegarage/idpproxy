package store

import (
	"context"
	"errors"
	"testing"
)

type fakeKMSStringAdapter struct {
	encryptBlob string
	encryptErr  error
	decryptText string
	decryptErr  error

	encryptCalled bool
	encryptPlain  string
	decryptCalled bool
	decryptBlob   string
}

func (f *fakeKMSStringAdapter) EncryptString(
	_ context.Context,
	plain string,
) (string, error) {
	f.encryptCalled = true
	f.encryptPlain = plain

	if f.encryptErr != nil {
		return "", f.encryptErr
	}

	return f.encryptBlob, nil
}

func (f *fakeKMSStringAdapter) DecryptString(
	_ context.Context,
	cipherB64 string,
) (string, error) {
	f.decryptCalled = true
	f.decryptBlob = cipherB64

	if f.decryptErr != nil {
		return "", f.decryptErr
	}

	return f.decryptText, nil
}

func TestNewKMSTokenEncryptor(t *testing.T) {
	t.Parallel()

	t.Run("creates_encryptor", func(t *testing.T) {
		t.Parallel()

		adapter := &fakeKMSStringAdapter{}

		got, err := NewKMSTokenEncryptor(
			"projects/test/locations/global/keyRings/ring/cryptoKeys/key",
			adapter,
		)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		if got.kid != "projects/test/locations/global/keyRings/ring/cryptoKeys/key" {
			t.Errorf("kid = %q, want configured key ID", got.kid)
		}
		if got.adapter != adapter {
			t.Error("adapter was not retained")
		}
	})

	t.Run("returns_error_for_empty_key_ID", func(t *testing.T) {
		t.Parallel()

		got, err := NewKMSTokenEncryptor("", &fakeKMSStringAdapter{})

		if !errors.Is(err, ErrEmptyKMSKeyID) {
			t.Errorf("error = %v, want %v", err, ErrEmptyKMSKeyID)
		}
		if got != nil {
			t.Errorf("encryptor = %#v, want nil", got)
		}
	})

	t.Run("returns_error_for_nil_adapter", func(t *testing.T) {
		t.Parallel()

		got, err := NewKMSTokenEncryptor(
			"projects/test/locations/global/keyRings/ring/cryptoKeys/key",
			nil,
		)

		if !errors.Is(err, ErrNilKMSStringAdapter) {
			t.Errorf("error = %v, want %v", err, ErrNilKMSStringAdapter)
		}
		if got != nil {
			t.Errorf("encryptor = %#v, want nil", got)
		}
	})
}

func TestKMSTokenEncryptor_EncryptString(t *testing.T) {
	t.Parallel()

	const kid = "projects/test/locations/global/keyRings/ring/cryptoKeys/key"

	t.Run("encrypts_and_returns_ciphertext", func(t *testing.T) {
		t.Parallel()

		adapter := &fakeKMSStringAdapter{
			encryptBlob: "base64-ciphertext",
		}
		encryptor, err := NewKMSTokenEncryptor(kid, adapter)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		got, err := encryptor.EncryptString(
			context.Background(),
			"github-access-token",
		)
		if err != nil {
			t.Fatalf("EncryptString() error = %v", err)
		}

		want := Ciphertext{
			KID:  kid,
			Blob: "base64-ciphertext",
		}
		if got != want {
			t.Errorf("EncryptString() = %#v, want %#v", got, want)
		}
		if !adapter.encryptCalled {
			t.Fatal("KMS EncryptString was not called")
		}
		if got := adapter.encryptPlain; got != "github-access-token" {
			t.Errorf("KMS EncryptString plain = %q, want %q", got, "github-access-token")
		}
	})

	t.Run("wraps_KMS_encryption_error", func(t *testing.T) {
		t.Parallel()

		kmsErr := errors.New("KMS unavailable")
		adapter := &fakeKMSStringAdapter{
			encryptErr: kmsErr,
		}
		encryptor, err := NewKMSTokenEncryptor(kid, adapter)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		got, err := encryptor.EncryptString(
			context.Background(),
			"github-access-token",
		)

		if !errors.Is(err, kmsErr) {
			t.Errorf("EncryptString() error = %v, want wrapped %v", err, kmsErr)
		}
		if got != (Ciphertext{}) {
			t.Errorf("EncryptString() = %#v, want empty Ciphertext", got)
		}
	})
}

func TestKMSTokenEncryptor_DecryptString(t *testing.T) {
	t.Parallel()

	const kid = "projects/test/locations/global/keyRings/ring/cryptoKeys/key"

	t.Run("decrypts_ciphertext", func(t *testing.T) {
		t.Parallel()

		adapter := &fakeKMSStringAdapter{
			decryptText: "github-access-token",
		}
		encryptor, err := NewKMSTokenEncryptor(kid, adapter)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		got, err := encryptor.DecryptString(context.Background(), Ciphertext{
			KID:  kid,
			Blob: "base64-ciphertext",
		})
		if err != nil {
			t.Fatalf("DecryptString() error = %v", err)
		}

		if got != "github-access-token" {
			t.Errorf("DecryptString() = %q, want %q", got, "github-access-token")
		}
		if !adapter.decryptCalled {
			t.Fatal("KMS DecryptString was not called")
		}
		if got := adapter.decryptBlob; got != "base64-ciphertext" {
			t.Errorf("KMS DecryptString blob = %q, want %q", got, "base64-ciphertext")
		}
	})

	t.Run("returns_error_when_key_ID_does_not_match", func(t *testing.T) {
		t.Parallel()

		adapter := &fakeKMSStringAdapter{}
		encryptor, err := NewKMSTokenEncryptor(kid, adapter)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		got, err := encryptor.DecryptString(context.Background(), Ciphertext{
			KID:  "projects/test/locations/global/keyRings/ring/cryptoKeys/other-key",
			Blob: "base64-ciphertext",
		})

		if !errors.Is(err, ErrKMSKeyIDMismatch) {
			t.Errorf("DecryptString() error = %v, want %v", err, ErrKMSKeyIDMismatch)
		}
		if got != "" {
			t.Errorf("DecryptString() = %q, want empty string", got)
		}
		if adapter.decryptCalled {
			t.Fatal("KMS DecryptString must not be called when key IDs differ")
		}
	})

	t.Run("returns_error_for_empty_ciphertext_blob", func(t *testing.T) {
		t.Parallel()

		adapter := &fakeKMSStringAdapter{}
		encryptor, err := NewKMSTokenEncryptor(kid, adapter)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		got, err := encryptor.DecryptString(context.Background(), Ciphertext{
			KID: kid,
		})

		if !errors.Is(err, ErrEmptyCiphertextBlob) {
			t.Errorf("DecryptString() error = %v, want %v", err, ErrEmptyCiphertextBlob)
		}
		if got != "" {
			t.Errorf("DecryptString() = %q, want empty string", got)
		}
		if adapter.decryptCalled {
			t.Fatal("KMS DecryptString must not be called for an empty blob")
		}
	})

	t.Run("wraps_KMS_decryption_error", func(t *testing.T) {
		t.Parallel()

		kmsErr := errors.New("KMS unavailable")
		adapter := &fakeKMSStringAdapter{
			decryptErr: kmsErr,
		}
		encryptor, err := NewKMSTokenEncryptor(kid, adapter)
		if err != nil {
			t.Fatalf("NewKMSTokenEncryptor() error = %v", err)
		}

		got, err := encryptor.DecryptString(context.Background(), Ciphertext{
			KID:  kid,
			Blob: "base64-ciphertext",
		})

		if !errors.Is(err, kmsErr) {
			t.Errorf("DecryptString() error = %v, want wrapped %v", err, kmsErr)
		}
		if got != "" {
			t.Errorf("DecryptString() = %q, want empty string", got)
		}
	})
}
