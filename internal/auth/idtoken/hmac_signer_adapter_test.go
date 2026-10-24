package idtoken

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/auth/signer"
)

func TestHMACSignerAdapter_SignJWT(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("signs payload and returns key ID", func(t *testing.T) {
		t.Parallel()

		const (
			key   = "test-signing-key"
			keyID = "test-key-id"
		)

		now := time.Now().UTC()
		hmacSigner := signer.NewHMACSigner([]byte(key), keyID)
		adapter := NewHMACSignerAdapter(hmacSigner)

		token, gotKeyID, err := adapter.SignJWT(ctx, map[string]any{
			"iss": "https://idpproxy.example.com",
			"sub": "user-123",
			"aud": "client-123",
			"iat": now.Unix(),
			"exp": now.Add(time.Hour).Unix(),
		})

		if err != nil {
			t.Fatalf("SignJWT() error = %v", err)
		}
		if token == "" {
			t.Fatal("SignJWT() returned an empty token")
		}
		if gotKeyID != keyID {
			t.Fatalf("kid = %q, want %q", gotKeyID, keyID)
		}

		result, err := hmacSigner.Verify(ctx, token, &signer.VerifyOptions{
			ExpectKID:  keyID,
			RequireTyp: true,
		})
		if err != nil {
			t.Fatalf("Verify() error = %v", err)
		}

		if got := result.Claims["iss"]; got != "https://idpproxy.example.com" {
			t.Fatalf("iss = %v, want %q", got, "https://idpproxy.example.com")
		}
		if got := result.Claims["sub"]; got != "user-123" {
			t.Fatalf("sub = %v, want %q", got, "user-123")
		}
		if got := result.Claims["aud"]; got != "client-123" {
			t.Fatalf("aud = %v, want %q", got, "client-123")
		}
	})

	t.Run("returns JSON marshal error", func(t *testing.T) {
		t.Parallel()

		hmacSigner := signer.NewHMACSigner(
			[]byte("test-signing-key"),
			"test-key-id",
		)
		adapter := NewHMACSignerAdapter(hmacSigner)

		_, _, err := adapter.SignJWT(ctx, map[string]any{
			"unsupported": make(chan struct{}),
		})

		if err == nil {
			t.Fatal("SignJWT() error = nil, want error")
		}
	})

	t.Run("returns signing error", func(t *testing.T) {
		t.Parallel()

		hmacSigner := signer.NewHMACSigner(nil, "test-key-id")
		adapter := NewHMACSignerAdapter(hmacSigner)

		_, _, err := adapter.SignJWT(ctx, map[string]any{
			"sub": "user-123",
		})

		if !errors.Is(err, signer.ErrEmptyKey) {
			t.Fatalf("error = %v, want %v", err, signer.ErrEmptyKey)
		}
	})
}

func TestHMACSignerAdapter_Algorithm(t *testing.T) {
	t.Parallel()

	hmacSigner := signer.NewHMACSigner(
		[]byte("test-signing-key"),
		"test-key-id",
	)
	adapter := NewHMACSignerAdapter(hmacSigner)

	if got := adapter.Algorithm(); got != signer.AlgHS256 {
		t.Fatalf("Algorithm() = %q, want %q", got, signer.AlgHS256)
	}
}
