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

		var k fernet.Key
		require.NoError(t, k.Generate())
		adapter, err := NewFernetAdapter(&k, 0)
		require.NoError(t, err)

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

		adapter, err := NewFernetAdapter(nil, 0)
		require.Error(t, err)
		require.Nil(t, adapter)
		require.ErrorIs(t, err, ErrNilKey)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		t.Parallel()

		var k fernet.Key
		require.NoError(t, k.Generate())
		adapter, err := NewFernetAdapter(&k, 0)
		require.NoError(t, err)

		_, err = adapter.DecryptString("not-a-valid-token")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired token")
	})

	t.Run("TTL", func(t *testing.T) {
		t.Parallel()

		var k fernet.Key
		require.NoError(t, k.Generate())
		adapter, err := NewFernetAdapter(&k, 1)
		require.NoError(t, err)

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
