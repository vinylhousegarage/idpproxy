package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type fakeConsumeStore struct {
	called   bool
	gotCode  string
	gotCID   string
	retCode  authcode.ProxyCode
	retError error
}

func (f *fakeConsumeStore) Save(
	context.Context,
	authcode.ProxyCode,
) error {
	panic("not used")
}

func (f *fakeConsumeStore) Consume(
	_ context.Context,
	proxyCode string,
	clientID string,
) (authcode.ProxyCode, error) {
	f.called = true
	f.gotCode = proxyCode
	f.gotCID = clientID

	return f.retCode, f.retError
}

func TestService_Consume(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success_returns_proxy_code", func(t *testing.T) {
		t.Parallel()

		want := authcode.ProxyCode{
			Code:      "code-abc",
			UserID:    "user-123",
			ClientID:  "client-xyz",
			ExpiresAt: time.Date(2026, time.October, 19, 1, 0, 0, 0, time.UTC),
		}
		store := &fakeConsumeStore{retCode: want}
		svc := &Service{store: store}

		got, err := svc.Consume(ctx, "code-abc", "client-xyz")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Fatalf("Consume() = %#v, want %#v", got, want)
		}

		if !store.called {
			t.Fatal("store.Consume was not called")
		}
		if store.gotCode != "code-abc" {
			t.Fatalf("code = %q, want %q", store.gotCode, "code-abc")
		}
		if store.gotCID != "client-xyz" {
			t.Fatalf("clientID = %q, want %q", store.gotCID, "client-xyz")
		}
	})

	t.Run("propagates_error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("invalid code")
		store := &fakeConsumeStore{retError: wantErr}
		svc := &Service{store: store}

		got, err := svc.Consume(ctx, "bad-code", "client-xyz")

		if err != wantErr {
			t.Fatalf("error = %v, want %v", err, wantErr)
		}
		if got != (authcode.ProxyCode{}) {
			t.Fatalf("Consume() = %#v, want zero ProxyCode", got)
		}
	})
}
