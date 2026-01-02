package service

import (
	"context"
	"testing"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type fakeStore struct {
	saved authcode.AuthCode
	err   error
}

func (f *fakeStore) Save(ctx context.Context, code authcode.AuthCode) error {
	f.saved = code
	return f.err
}

func (f *fakeStore) Consume(ctx context.Context, code string, clientID string) (string, error) {
	return "", nil
}

func TestService_Issue(t *testing.T) {
	t.Parallel()

	t.Run("successfully issues auth code and saves it", func(t *testing.T) {
		t.Parallel()

		fs := &fakeStore{}
		svc := &Service{
			store: fs,
		}

		err := svc.Issue(context.Background(), "user-1", "client-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if fs.saved.UserID != "user-1" {
			t.Errorf("UserID mismatch: got=%s", fs.saved.UserID)
		}

		if fs.saved.ClientID != "client-1" {
			t.Errorf("ClientID mismatch: got=%s", fs.saved.ClientID)
		}

		if fs.saved.Code == "" {
			t.Errorf("Code should not be empty")
		}

		if time.Until(fs.saved.ExpiresAt) <= 0 {
			t.Errorf("ExpiresAt should be in the future: got=%v", fs.saved.ExpiresAt)
		}
	})

	t.Run("returns error when store.Save fails", func(t *testing.T) {
		t.Parallel()

		expectedErr := context.DeadlineExceeded

		fs := &fakeStore{
			err: expectedErr,
		}
		svc := &Service{
			store: fs,
		}

		err := svc.Issue(context.Background(), "user-1", "client-1")
		if err != expectedErr {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})
}
