package encrypt

import (
	"testing"
	"time"

	"github.com/fernet/fernet-go"
	"github.com/stretchr/testify/require"
)

func TestFernetAdapter(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		keys := fernet.MustGenerateKeys()
		key := keys[0]
		adapter := NewFernetAdapter(key, 0)

		plain := "hello world"
		token, err := adapter.EncryptString(plain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		out, err := adapter.DecryptString(token)
		require.NoError(t, err)
		require.Equal(t, plain, out)
	})

	t.Run("KeyNil", func(t *testing.T) {
		t.Parallel()

		adapter := NewFernetAdapter(nil, 0)

		_, err := adapter.EncryptString("test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "fernet key is nil")

		_, err = adapter.DecryptString("dummy")
		require.Error(t, err)
		require.Contains(t, err.Error(), "fernet key is nil")
	})

	t.Run("InvalidToken", func(t *testing.T) {
		t.Parallel()

		keys := fernet.MustGenerateKeys()
		key := keys[0]
		adapter := NewFernetAdapter(key, 0)

		_, err := adapter.DecryptString("not-a-valid-token")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired token")
	})

	t.Run("TTL", func(t *testing.T) {
		t.Parallel()

		keys := fernet.MustGenerateKeys()
		key := keys[0]
		adapter := NewFernetAdapter(key, 1)

		plain := "short lived"
		token, err := adapter.EncryptString(plain)
		require.NoError(t, err)

		out, err := adapter.DecryptString(token)
		require.NoError(t, err)
		require.Equal(t, plain, out)

		time.Sleep(1100 * time.Millisecond)
		_, err = adapter.DecryptString(token)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired token")
	})
}
