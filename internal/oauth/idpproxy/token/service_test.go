package token

import (
	"context"
	"strings"
	"testing"
	"time"
)

type fixedClock struct {
	t time.Time
}

func (f fixedClock) Now() time.Time {
	return f.t
}

type mockStore struct {
	code *AuthCode
	err  error
}

func (m *mockStore) Consume(ctx context.Context, code, clientID string) (*AuthCode, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.code, nil
}

func newTestService() *Service {
	return &Service{
		Store: &mockStore{err: ErrInvalidGrant},
		Clock: fixedClock{t: time.Now()},
	}
}

func newTestServiceWithExpiredCode() *Service {
	return &Service{
		Store: &mockStore{
			code: &AuthCode{
				UserID:    "user1",
				ClientID:  "client-1",
				ExpiresAt: time.Now().Add(-time.Hour),
			},
		},
		Clock: fixedClock{t: time.Now()},
	}
}

func newTestServiceWithValidCode() *Service {
	return &Service{
		Store: &mockStore{
			code: &AuthCode{
				UserID:    "user1",
				ClientID:  "client-1",
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
		Clock: fixedClock{t: time.Now()},
	}
}

func TestService_Exchange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("unsupported grant_type returns ErrUnsupportedGrantType", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()

		_, err := svc.Exchange(ctx, TokenRequest{
			GrantType: "password",
		})

		if err != ErrUnsupportedGrantType {
			t.Fatalf("expected ErrUnsupportedGrantType, got %v", err)
		}
	})

	t.Run("invalid auth code returns ErrInvalidGrant", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()

		_, err := svc.Exchange(ctx, TokenRequest{
			GrantType: "authorization_code",
			Code:      "no-such-code",
			ClientID:  "client-1",
		})

		if err != ErrInvalidGrant {
			t.Fatalf("expected ErrInvalidGrant, got %v", err)
		}
	})

	t.Run("expired auth code returns ErrInvalidGrant", func(t *testing.T) {
		t.Parallel()

		svc := newTestServiceWithExpiredCode()

		_, err := svc.Exchange(ctx, TokenRequest{
			GrantType: "authorization_code",
			Code:      "expired-code",
			ClientID:  "client-1",
		})

		if err != ErrInvalidGrant {
			t.Fatalf("expected ErrInvalidGrant, got %v", err)
		}
	})

	t.Run("valid request returns id_token", func(t *testing.T) {
		t.Parallel()

		svc := newTestServiceWithValidCode()

		resp, err := svc.Exchange(ctx, TokenRequest{
			GrantType: "authorization_code",
			Code:      "valid-code",
			ClientID:  "client-1",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp == nil {
			t.Fatal("response should not be nil")
		}

		if resp.IDToken == "" {
			t.Fatal("id_token should not be empty")
		}

		if !strings.Contains(resp.IDToken, "user1") {
			t.Fatalf("id_token should contain user info")
		}
	})
}
