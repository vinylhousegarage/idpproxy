package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type fakeProxyCodeConsumer struct {
	called      bool
	gotCode     string
	gotClientID string
	result      authcode.ProxyCode
	err         error
}

func (f *fakeProxyCodeConsumer) Consume(
	_ context.Context,
	code string,
	clientID string,
) (authcode.ProxyCode, error) {
	f.called = true
	f.gotCode = code
	f.gotClientID = clientID

	return f.result, f.err
}

func TestAuthCodeServiceAdapter_Consume(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("converts_proxy_code_to_auth_code", func(t *testing.T) {
		t.Parallel()

		expiresAt := time.Date(
			2026,
			time.October,
			20,
			1,
			0,
			0,
			0,
			time.UTC,
		)
		consumer := &fakeProxyCodeConsumer{
			result: authcode.ProxyCode{
				Code:      "proxy-code-123",
				UserID:    "user-123",
				ClientID:  "client-123",
				ExpiresAt: expiresAt,
			},
		}
		adapter := NewAuthCodeServiceAdapter(consumer)

		got, err := adapter.Consume(ctx, "proxy-code-123", "client-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !consumer.called {
			t.Fatal("Consume was not called")
		}
		if consumer.gotCode != "proxy-code-123" {
			t.Fatalf("code = %q, want %q", consumer.gotCode, "proxy-code-123")
		}
		if consumer.gotClientID != "client-123" {
			t.Fatalf(
				"clientID = %q, want %q",
				consumer.gotClientID,
				"client-123",
			)
		}

		want := &AuthCode{
			UserID:    "user-123",
			ClientID:  "client-123",
			ExpiresAt: expiresAt,
		}
		if *got != *want {
			t.Fatalf("Consume() = %#v, want %#v", got, want)
		}
	})

	t.Run("propagates_consume_error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("proxy code not found")
		consumer := &fakeProxyCodeConsumer{
			err: wantErr,
		}
		adapter := NewAuthCodeServiceAdapter(consumer)

		got, err := adapter.Consume(ctx, "missing-code", "client-123")

		if err != wantErr {
			t.Fatalf("error = %v, want %v", err, wantErr)
		}
		if got != nil {
			t.Fatalf("Consume() = %#v, want nil", got)
		}
	})
}
