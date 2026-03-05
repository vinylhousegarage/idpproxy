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

	t.Run("returns ErrExpired and deletes code when auth code is expired", func(t *testing.T) {
		t.Parallel()

		s := NewMemoryStore()
		pc := authcode.ProxyCode{
			Code:      "code-expired",
			UserID:    "user-1",
			ClientID:  "client-1",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}

		_ = s.Save(ctx, pc)

		_, err := s.Consume(ctx, "code-expired", "client-1")
		if err != ErrExpired {
			t.Fatalf("expected ErrExpired, got %v", err)
		}

		_, err = s.Consume(ctx, "code-expired", "client-1")
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound after expiration delete, got %v", err)
		}
	})

	t.Run("returns ErrClientMismatch and does not delete code when client id does not match", func(t *testing.T) {
		t.Parallel()

		s := NewMemoryStore()
		pc := authcode.ProxyCode{
			Code:      "code-client",
			UserID:    "user-1",
			ClientID:  "client-1",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		_ = s.Save(ctx, pc)

		_, err := s.Consume(ctx, "code-client", "client-2")
		if err != ErrClientMismatch {
			t.Fatalf("expected ErrClientMismatch, got %v", err)
		}

		uid, err := s.Consume(ctx, "code-client", "client-1")
		if err != nil {
			t.Fatalf("expected success after mismatch, got %v", err)
		}
		if uid != "user-1" {
			t.Fatalf("unexpected user id: %s", uid)
		}
	})

	t.Run("deletes code after successful consume", func(t *testing.T) {
		t.Parallel()

		s := NewMemoryStore()
		pc := authcode.ProxyCode{
			Code:      "code-ok",
			UserID:    "user-1",
			ClientID:  "client-1",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		_ = s.Save(ctx, pc)

		uid, err := s.Consume(ctx, "code-ok", "client-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if uid != "user-1" {
			t.Fatalf("unexpected user id: %s", uid)
		}

		_, err = s.Consume(ctx, "code-ok", "client-1")
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound after consume, got %v", err)
		}
	})
}
