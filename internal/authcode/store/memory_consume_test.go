package store

import (
	"context"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

func TestMemoryStore_Consume(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("returns ErrNotFound when code does not exist", func(t *testing.T) {
		t.Parallel()

		s := NewMemoryStore()

		_, err := s.Consume(ctx, "no-such-code", "client-1")
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrExpired when auth code is expired", func(t *testing.T) {
		t.Parallel()

		s := NewMemoryStore()
		code := authcode.AuthCode{
			Code:      "code-expired",
			UserID:    "user-1",
			ClientID:  "client-1",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}

		_ = s.Save(ctx, code)

		_, err := s.Consume(ctx, "code-expired", "client-1")
		if err != ErrExpired {
			t.Fatalf("expected ErrExpired, got %v", err)
		}
	})

	t.Run("returns ErrClientMismatch when client id does not match", func(t *testing.T) {
		t.Parallel()

		s := NewMemoryStore()
		code := authcode.AuthCode{
			Code:      "code-client",
			UserID:    "user-1",
			ClientID:  "client-1",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		_ = s.Save(ctx, code)

		_, err := s.Consume(ctx, "code-client", "client-2")
		if err != ErrClientMismatch {
			t.Fatalf("expected ErrClientMismatch, got %v", err)
		}
	})
}
