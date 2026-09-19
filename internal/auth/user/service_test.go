package user

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/auth"
)

type fakeFirebaseAuthClient struct {
	getUserRecord  *auth.UserRecord
	getUserErr     error
	createRecord   *auth.UserRecord
	createErr      error
	updateRecord   *auth.UserRecord
	updateErr      error
	getCalled      bool
	getProviderID  string
	getProviderUID string
	createCalled   bool
	updateCalled   bool
	updateUID      string
}

func (f *fakeFirebaseAuthClient) GetUserByProviderUID(
	_ context.Context,
	providerID string,
	providerUID string,
) (*auth.UserRecord, error) {
	f.getCalled = true
	f.getProviderID = providerID
	f.getProviderUID = providerUID

	if f.getUserErr != nil {
		return nil, f.getUserErr
	}

	return f.getUserRecord, nil
}

func (f *fakeFirebaseAuthClient) CreateUser(
	_ context.Context,
	_ *auth.UserToCreate,
) (*auth.UserRecord, error) {
	f.createCalled = true

	if f.createErr != nil {
		return nil, f.createErr
	}

	return f.createRecord, nil
}

func (f *fakeFirebaseAuthClient) UpdateUser(
	_ context.Context,
	uid string,
	_ *auth.UserToUpdate,
) (*auth.UserRecord, error) {
	f.updateCalled = true
	f.updateUID = uid

	if f.updateErr != nil {
		return nil, f.updateErr
	}

	return f.updateRecord, nil
}

func firebaseUserRecord(uid string) *auth.UserRecord {
	return &auth.UserRecord{
		UserInfo: &auth.UserInfo{
			UID: uid,
		},
	}
}

func TestService_UpsertFromGitHub(t *testing.T) {
	t.Parallel()

	const (
		githubID  = int64(123456)
		githubUID = "123456"
		login     = "vinylhousegarage"
		email     = "user@example.com"
	)

	t.Run("updates_existing_GitHub_user_and_returns_Firebase_UID", func(t *testing.T) {
		t.Parallel()

		client := &fakeFirebaseAuthClient{
			getUserRecord: firebaseUserRecord("firebase-user-123"),
			updateRecord:  firebaseUserRecord("firebase-user-123"),
		}
		service := NewService(client)

		got, err := service.UpsertFromGitHub(
			context.Background(),
			githubID,
			login,
			email,
		)
		if err != nil {
			t.Fatalf("UpsertFromGitHub() error = %v", err)
		}

		if got != "firebase-user-123" {
			t.Errorf("UpsertFromGitHub() = %q, want %q", got, "firebase-user-123")
		}
		if !client.getCalled {
			t.Fatal("GetUserByProviderUID was not called")
		}
		if got := client.getProviderID; got != "github.com" {
			t.Errorf("provider ID = %q, want %q", got, "github.com")
		}
		if got := client.getProviderUID; got != githubUID {
			t.Errorf("provider UID = %q, want %q", got, githubUID)
		}
		if client.createCalled {
			t.Fatal("CreateUser must not be called for an existing user")
		}
		if !client.updateCalled {
			t.Fatal("UpdateUser was not called")
		}
		if got := client.updateUID; got != "firebase-user-123" {
			t.Errorf("UpdateUser UID = %q, want %q", got, "firebase-user-123")
		}
	})

	t.Run("creates_and_links_new_GitHub_user", func(t *testing.T) {
		t.Parallel()

		notFoundErr := errors.New("user not found")
		client := &fakeFirebaseAuthClient{
			getUserErr:   notFoundErr,
			createRecord: firebaseUserRecord("firebase-user-456"),
			updateRecord: firebaseUserRecord("firebase-user-456"),
		}
		service := NewService(client)
		service.isUserNotFound = func(err error) bool {
			return errors.Is(err, notFoundErr)
		}

		got, err := service.UpsertFromGitHub(
			context.Background(),
			githubID,
			login,
			email,
		)
		if err != nil {
			t.Fatalf("UpsertFromGitHub() error = %v", err)
		}

		if got != "firebase-user-456" {
			t.Errorf("UpsertFromGitHub() = %q, want %q", got, "firebase-user-456")
		}
		if !client.createCalled {
			t.Fatal("CreateUser was not called")
		}
		if !client.updateCalled {
			t.Fatal("UpdateUser was not called to link the GitHub provider")
		}
		if got := client.updateUID; got != "firebase-user-456" {
			t.Errorf("UpdateUser UID = %q, want %q", got, "firebase-user-456")
		}
	})

	t.Run("returns_get_user_error", func(t *testing.T) {
		t.Parallel()

		getErr := errors.New("Firebase unavailable")
		client := &fakeFirebaseAuthClient{
			getUserErr: getErr,
		}
		service := NewService(client)
		service.isUserNotFound = func(error) bool {
			return false
		}

		got, err := service.UpsertFromGitHub(
			context.Background(),
			githubID,
			login,
			email,
		)

		if !errors.Is(err, getErr) {
			t.Errorf("UpsertFromGitHub() error = %v, want wrapped %v", err, getErr)
		}
		if got != "" {
			t.Errorf("UpsertFromGitHub() = %q, want empty string", got)
		}
		if client.createCalled {
			t.Fatal("CreateUser must not be called after a get-user failure")
		}
		if client.updateCalled {
			t.Fatal("UpdateUser must not be called after a get-user failure")
		}
	})

	t.Run("returns_create_user_error", func(t *testing.T) {
		t.Parallel()

		notFoundErr := errors.New("user not found")
		createErr := errors.New("Firebase unavailable")
		client := &fakeFirebaseAuthClient{
			getUserErr: notFoundErr,
			createErr:  createErr,
		}
		service := NewService(client)
		service.isUserNotFound = func(err error) bool {
			return errors.Is(err, notFoundErr)
		}

		got, err := service.UpsertFromGitHub(
			context.Background(),
			githubID,
			login,
			email,
		)

		if !errors.Is(err, createErr) {
			t.Errorf("UpsertFromGitHub() error = %v, want wrapped %v", err, createErr)
		}
		if got != "" {
			t.Errorf("UpsertFromGitHub() = %q, want empty string", got)
		}
		if !client.createCalled {
			t.Fatal("CreateUser was not called")
		}
		if client.updateCalled {
			t.Fatal("UpdateUser must not be called when CreateUser fails")
		}
	})

	t.Run("returns_update_user_error", func(t *testing.T) {
		t.Parallel()

		updateErr := errors.New("Firebase unavailable")
		client := &fakeFirebaseAuthClient{
			getUserRecord: firebaseUserRecord("firebase-user-123"),
			updateErr:     updateErr,
		}
		service := NewService(client)

		got, err := service.UpsertFromGitHub(
			context.Background(),
			githubID,
			login,
			email,
		)

		if !errors.Is(err, updateErr) {
			t.Errorf("UpsertFromGitHub() error = %v, want wrapped %v", err, updateErr)
		}
		if got != "" {
			t.Errorf("UpsertFromGitHub() = %q, want empty string", got)
		}
		if !client.updateCalled {
			t.Fatal("UpdateUser was not called")
		}
	})
}
