package client

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateClientID(t *testing.T) {
	t.Run("returns a base64url encoded client ID", func(t *testing.T) {
		clientID, err := GenerateClientID()

		require.NoError(t, err)
		require.NotEmpty(t, clientID)

		decoded, err := base64.RawURLEncoding.DecodeString(clientID)
		require.NoError(t, err)
		require.Len(t, decoded, 32)
	})

	t.Run("generates different client IDs", func(t *testing.T) {
		first, err := GenerateClientID()
		require.NoError(t, err)

		second, err := GenerateClientID()
		require.NoError(t, err)

		require.NotEqual(t, first, second)
	})
}
