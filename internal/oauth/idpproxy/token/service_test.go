package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/auth/idtoken"
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

func (m *mockStore) Consume(
	_ context.Context,
	_ string,
	_ string,
) (*AuthCode, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.code, nil
}

type mockIDTokenIssuer struct {
	called bool
	input  *idtoken.IDTokenInput
	token  string
	err    error
}

func (m *mockIDTokenIssuer) Issue(
	_ context.Context,
	in *idtoken.IDTokenInput,
) (string, string, error) {
	m.called = true
	m.input = in

	if m.err != nil {
		return "", "", m.err
	}

	return m.token, "test-key-id", nil
}

func newTestService() *Service {
	return &Service{
		Store: &mockStore{err: ErrInvalidGrant},
		Clock: fixedClock{t: time.Now()},
		IDTokenIssuer: &mockIDTokenIssuer{
			token: "signed-id-token",
		},
		IDTokenTTL: time.Hour,
	}
}

func TestService_Exchange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Date(
		2026,
		time.October,
		23,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	validRequest := TokenRequest{
		GrantType: "authorization_code",
		Code:      "valid-code",
		ClientID:  "client-1",
	}

	t.Run("unsupported grant type returns ErrUnsupportedGrantType", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()

		_, err := svc.Exchange(ctx, TokenRequest{
			GrantType: "password",
		})

		if !errors.Is(err, ErrUnsupportedGrantType) {
			t.Fatalf(
				"error = %v, want %v",
				err,
				ErrUnsupportedGrantType,
			)
		}
	})

	t.Run("invalid auth code returns ErrInvalidGrant", func(t *testing.T) {
		t.Parallel()

		svc := newTestService()

		_, err := svc.Exchange(ctx, validRequest)

		if !errors.Is(err, ErrInvalidGrant) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidGrant)
		}
	})

	t.Run("expired auth code returns ErrInvalidGrant", func(t *testing.T) {
		t.Parallel()

		issuer := &mockIDTokenIssuer{
			token: "signed-id-token",
		}
		svc := &Service{
			Store: &mockStore{
				code: &AuthCode{
					UserID:    "user1",
					ClientID:  "client-1",
					ExpiresAt: now.Add(-time.Hour),
				},
			},
			Clock:         fixedClock{t: now},
			IDTokenIssuer: issuer,
			IDTokenTTL:    time.Hour,
		}

		_, err := svc.Exchange(ctx, validRequest)

		if !errors.Is(err, ErrInvalidGrant) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidGrant)
		}
		if issuer.called {
			t.Fatal("ID token issuer must not be called for an expired code")
		}
	})

	t.Run("auth code expiring now returns ErrInvalidGrant", func(t *testing.T) {
		t.Parallel()

		issuer := &mockIDTokenIssuer{
			token: "signed-id-token",
		}
		svc := &Service{
			Store: &mockStore{
				code: &AuthCode{
					UserID:    "user1",
					ClientID:  "client-1",
					ExpiresAt: now,
				},
			},
			Clock:         fixedClock{t: now},
			IDTokenIssuer: issuer,
			IDTokenTTL:    time.Hour,
		}

		_, err := svc.Exchange(ctx, validRequest)

		if !errors.Is(err, ErrInvalidGrant) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidGrant)
		}
		if issuer.called {
			t.Fatal("ID token issuer must not be called for an expired code")
		}
	})
}
