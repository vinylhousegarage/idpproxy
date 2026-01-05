package service

import (
	"context"
	"errors"
	"testing"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type fakeConsumeStore struct {
	called   bool
	gotCode  string
	gotCID   string
	retUID   string
	retError error
}

func (f *fakeConsumeStore) Save(
	ctx context.Context,
	code authcode.AuthCode,
) error {
	panic("not used")
}

func (f *fakeConsumeStore) Consume(
	ctx context.Context,
	code string,
	clientID string,
) (string, error) {
	f.called = true
	f.gotCode = code
	f.gotCID = clientID
	return f.retUID, f.retError
}

func TestService_Consume(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success_returns_user_id", func(t *testing.T) {
		t.Parallel()

		store := &fakeConsumeStore{
			retUID:   "user-123",
			retError: nil,
		}
		svc := &Service{store: store}

		uid, err := svc.Consume(ctx, "code-abc", "client-xyz")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if uid != "user-123" {
			t.Fatalf("unexpected uid: got=%s", uid)
		}

		if !store.called {
			t.Fatalf("store.Consume was not called")
		}
		if store.gotCode != "code-abc" {
			t.Fatalf("unexpected code: got=%s", store.gotCode)
		}
		if store.gotCID != "client-xyz" {
			t.Fatalf("unexpected clientID: got=%s", store.gotCID)
		}
	})

	t.Run("error_is_propagated", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("invalid code")
		store := &fakeConsumeStore{
			retError: wantErr,
		}
		svc := &Service{store: store}

		uid, err := svc.Consume(ctx, "bad-code", "client-xyz")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err != wantErr {
			t.Fatalf("unexpected error: %v", err)
		}
		if uid != "" {
			t.Fatalf("unexpected uid: got=%s", uid)
		}
	})
}
