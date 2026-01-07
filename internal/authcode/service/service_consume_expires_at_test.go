package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type fakeStore struct {
	authCode authcode.AuthCode
	err      error
}

func (f *fakeStore) Get(
	_ context.Context,
	_ string,
	_ string,
) (authcode.AuthCode, error) {
	if f.err != nil {
		return authcode.AuthCode{}, f.err
	}
	return f.authCode, nil
}

func (f *fakeStore) Save(context.Context, authcode.AuthCode) error {
	return nil
}

func (f *fakeStore) Consume(context.Context, string, string) (string, error) {
	return "", nil
}

func TestService_Consume_ExpiresAt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("returns ErrExpired when auth code is expired", func(t *testing.T) {
		t.Parallel()

		store := &fakeStore{
			authCode: authcode.AuthCode{
				Code:      "code-expired",
				UserID:    "user-1",
				ClientID:  "client-1",
				ExpiresAt: now.Add(-time.Second),
			},
		}

		svc := &Service{
			store: store,
			now:   func() time.Time { return now },
		}

		_, err := svc.Consume(ctx, "code-expired", "client-1")
		if !errors.Is(err, ErrExpired) {
			t.Fatalf("expected ErrExpired, got %v", err)
		}
	})

	t.Run("returns userID when auth code is not expired", func(t *testing.T) {
		t.Parallel()

		store := &fakeStore{
			authCode: authcode.AuthCode{
				Code:      "code-valid",
				UserID:    "user-1",
				ClientID:  "client-1",
				ExpiresAt: now.Add(time.Minute),
			},
		}

		svc := &Service{
			store: store,
			now:   func() time.Time { return now },
		}

		uid, err := svc.Consume(ctx, "code-valid", "client-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if uid != "user-1" {
			t.Fatalf("unexpected userID: got=%s", uid)
		}
	})
}
