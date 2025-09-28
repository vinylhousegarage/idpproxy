package signer

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestHMACSigner_Sign(t *testing.T) {
	t.Parallel()

	t.Run("success-empty-payload", func(t *testing.T) {
		t.Parallel()

		s := NewHMACSigner([]byte(testKey), testKid123)
		token, kid, err := s.Sign(context.Background(), nil)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.Equal(t, testKid123, kid)
	})

	t.Run("success-with-payload", func(t *testing.T) {
		t.Parallel()

		s := NewHMACSigner([]byte(testKey), testKidXYZ)
		payload := []byte(`{"sub":"user-1"}`)
		token, kid, err := s.Sign(context.Background(), payload)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.Equal(t, testKidXYZ, kid)
	})

	t.Run("error-empty-key", func(t *testing.T) {
		t.Parallel()

		s := NewHMACSigner(nil, testKidErr)
		_, _, err := s.Sign(context.Background(), nil)
		require.ErrorIs(t, err, ErrEmptyKey)
	})

	t.Run("error-invalid-payload", func(t *testing.T) {
		t.Parallel()

		s := NewHMACSigner([]byte(testKey), testKidErr2)
		_, _, err := s.Sign(context.Background(), []byte("{invalid-json"))
		require.ErrorIs(t, err, ErrInvalidPayload)
	})

	t.Run("header-kid-absent-when-empty", func(t *testing.T) {
		t.Parallel()

		s := NewHMACSigner([]byte(testKey), "")
		token, _, err := s.Sign(context.Background(), nil)
		require.NoError(t, err)

		parsed, err := jwt.Parse(token, func(tk *jwt.Token) (any, error) { return []byte(testKey), nil })
		require.NoError(t, err)
		require.True(t, parsed.Valid)
		_, hasKid := parsed.Header["kid"]
		require.False(t, hasKid, "kid header should be absent when keyID is empty")
	})

	t.Run("claims-use-fixed-now", func(t *testing.T) {
		t.Parallel()

		fixed := time.Now().Add(time.Hour)
		s := NewHMACSigner([]byte(testKey), testKidXYZ)
		s.now = fixedNow(fixed)

		tok, _, err := s.Sign(context.Background(), []byte(`{"sub":"u1"}`))
		require.NoError(t, err)

		parsed, err := jwt.Parse(tok, func(tk *jwt.Token) (any, error) { return []byte(testKey), nil })
		require.NoError(t, err)
		require.True(t, parsed.Valid)

		claims := parsed.Claims.(jwt.MapClaims)
		require.EqualValues(t, fixed.Unix(), claims["iat"])
		require.EqualValues(t, fixed.Add(24*time.Hour).Unix(), claims["exp"])
		require.Equal(t, "u1", claims["sub"])
		require.Equal(t, "JWT", parsed.Header["typ"])
		require.Equal(t, testKidXYZ, parsed.Header["kid"])
	})
}
