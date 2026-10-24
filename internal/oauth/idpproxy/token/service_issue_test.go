package token

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestService_Exchange_IDToken(t *testing.T) {
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

	t.Run("valid request issues and returns id token", func(t *testing.T) {
		t.Parallel()

		issuer := &mockIDTokenIssuer{
			token: "signed-id-token",
		}
		svc := &Service{
			Store: &mockStore{
				code: &AuthCode{
					UserID:    "user1",
					ClientID:  "client-1",
					ExpiresAt: now.Add(time.Hour),
				},
			},
			Clock:         fixedClock{t: now},
			IDTokenIssuer: issuer,
			IDTokenTTL:    30 * time.Minute,
		}

		resp, err := svc.Exchange(ctx, validRequest)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("response should not be nil")
		}
		if resp.IDToken != "signed-id-token" {
			t.Fatalf(
				"id_token = %q, want %q",
				resp.IDToken,
				"signed-id-token",
			)
		}
		if !issuer.called {
			t.Fatal("ID token issuer was not called")
		}
		if issuer.input == nil {
			t.Fatal("ID token issuer input is nil")
		}

		if issuer.input.UserID != "user1" {
			t.Fatalf(
				"UserID = %q, want %q",
				issuer.input.UserID,
				"user1",
			)
		}
		if issuer.input.ClientID != "client-1" {
			t.Fatalf(
				"ClientID = %q, want %q",
				issuer.input.ClientID,
				"client-1",
			)
		}
		if !issuer.input.Now.Equal(now) {
			t.Fatalf("Now = %v, want %v", issuer.input.Now, now)
		}
		if issuer.input.TTL != 30*time.Minute {
			t.Fatalf(
				"TTL = %v, want %v",
				issuer.input.TTL,
				30*time.Minute,
			)
		}
	})

	t.Run("id token issue error is propagated", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("signing failed")
		issuer := &mockIDTokenIssuer{
			err: wantErr,
		}
		svc := &Service{
			Store: &mockStore{
				code: &AuthCode{
					UserID:    "user1",
					ClientID:  "client-1",
					ExpiresAt: now.Add(time.Hour),
				},
			},
			Clock:         fixedClock{t: now},
			IDTokenIssuer: issuer,
			IDTokenTTL:    time.Hour,
		}

		resp, err := svc.Exchange(ctx, validRequest)

		if !errors.Is(err, wantErr) {
			t.Fatalf("error = %v, want %v", err, wantErr)
		}
		if resp != nil {
			t.Fatalf("response = %#v, want nil", resp)
		}
		if !issuer.called {
			t.Fatal("ID token issuer was not called")
		}
	})
}
