package signer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testKey    = "secret"
	testKid123 = "kid-123"
)

func TestHMACSigner_InfoMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  []byte
		kid  string
	}{
		{"basic", []byte("k1"), "kid-1"},
		{"empty-kid", []byte("k2"), ""},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewHMACSigner(tt.key, tt.kid)
			require.Equal(t, AlgHS256, s.Alg())
			require.Equal(t, tt.kid, s.KeyID())
		})
	}

	t.Run("key-cloned-defensively", func(t *testing.T) {
		t.Parallel()

		orig := []byte(testKey)
		s := NewHMACSigner(orig, testKid123)

		fixed := time.Unix(1_700_000_000, 0)
		s.now = func() time.Time { return fixed }

		payload := []byte(`{"sub":"u1"}`)
		tok1, _, err := s.Sign(context.Background(), payload)
		require.NoError(t, err)

		orig[0] = 'X'

		tok2, _, err := s.Sign(context.Background(), payload)
		require.NoError(t, err)
		require.Equal(t, tok1, tok2, "signer must hold a cloned key buffer")
	})
}
