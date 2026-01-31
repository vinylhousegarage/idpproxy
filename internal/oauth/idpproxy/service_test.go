package token

import (
	"context"
	"testing"
)

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

	t.Run("expired auth code returns ErrExpiredCode", func(t *testing.T) {
		t.Parallel()

		svc := newTestServiceWithExpiredCode()

		_, err := svc.Exchange(ctx, TokenRequest{
			GrantType: "authorization_code",
			Code:      "expired-code",
			ClientID:  "client-1",
		})

		if err != ErrExpiredCode {
			t.Fatalf("expected ErrExpiredCode, got %v", err)
		}
	})

	t.Run("invalid client returns ErrInvalidClient", func(t *testing.T) {
		t.Parallel()

		svc := newTestServiceWithValidCode()

		_, err := svc.Exchange(ctx, TokenRequest{
			GrantType:    "authorization_code",
			Code:         "valid-code",
			ClientID:     "client-1",
			ClientSecret: "wrong-secret",
		})

		if err != ErrInvalidClient {
			t.Fatalf("expected ErrInvalidClient, got %v", err)
		}
	})

	t.Run("valid request returns token response", func(t *testing.T) {
		t.Parallel()

		svc := newTestServiceWithValidCode()

		resp, err := svc.Exchange(ctx, TokenRequest{
			GrantType:    "authorization_code",
			Code:         "valid-code",
			ClientID:     "client-1",
			ClientSecret: "secret",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.AccessToken == "" {
			t.Fatal("access token should not be empty")
		}
	})
}
