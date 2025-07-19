package fernet

import (
	"testing"

	"github.com/fernet/fernet-go"
	"github.com/stretchr/testify/require"
)

func TestDecodeFernetKey(t *testing.T) {
	t.Parallel()

	t.Run("valid key", func(t *testing.T) {
		t.Parallel()

		var key fernet.Key
		err := key.Generate()
		require.NoError(t, err)

		encoded := key.Encode()
		decoded, err := DecodeFernetKey(encoded)
		require.NoError(t, err)
		require.Equal(t, key, *decoded)
	})

	t.Run("invalid base64", func(t *testing.T) {
		t.Parallel()

		_, err := DecodeFernetKey("invalid_base64_token")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to decode")
	})
}
