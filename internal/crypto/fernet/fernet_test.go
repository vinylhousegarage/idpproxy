package fernet

import (
	"testing"

	"github.com/fernet/fernet-go"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	key := &fernet.Key{}
	err := key.Generate()
	require.NoError(t, err)

	plaintext := []byte("secret data")
	token, err := Encrypt(key, plaintext)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	decrypted, err := Decrypt(key, token)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestDecrypt_Errors(t *testing.T) {
	t.Parallel()

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()

		key := &fernet.Key{}
		require.NoError(t, key.Generate())

		invalidToken := "this_is_not_a_valid_token"
		_, err := Decrypt(key, invalidToken)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired token")
	})

	t.Run("empty key", func(t *testing.T) {
		t.Parallel()

		var emptyKey *fernet.Key
		token := "dummy"
		_, err := Decrypt(emptyKey, token)
		require.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		t.Parallel()

		key := &fernet.Key{}
		require.NoError(t, key.Generate())

		plaintext := []byte("secret")
		token := fernet.EncryptAndSign(plaintext, key)

		_, err := Decrypt(key, string(token))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired token")
	})
}
