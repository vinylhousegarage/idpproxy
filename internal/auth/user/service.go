package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"firebase.google.com/go/v4/auth"
)

const githubProviderID = "github.com"

var (
	ErrNilFirebaseAuthClient  = errors.New("nil Firebase Auth client")
	ErrInvalidGitHubID        = errors.New("invalid GitHub ID")
	ErrEmptyGitHubLogin       = errors.New("empty GitHub login")
	ErrFirebaseUserMissingUID = errors.New("Firebase user has no UID")
)

type firebaseAuthClient interface {
	GetUserByProviderUID(
		ctx context.Context,
		providerID string,
		providerUID string,
	) (*auth.UserRecord, error)

	CreateUser(
		ctx context.Context,
		user *auth.UserToCreate,
	) (*auth.UserRecord, error)

	UpdateUser(
		ctx context.Context,
		uid string,
		user *auth.UserToUpdate,
	) (*auth.UserRecord, error)
}

type Service struct {
	client         firebaseAuthClient
	isUserNotFound func(error) bool
}

func NewService(client firebaseAuthClient) *Service {
	return &Service{
		client:         client,
		isUserNotFound: auth.IsUserNotFound,
	}
}

func (s *Service) UpsertFromGitHub(
	ctx context.Context,
	githubID int64,
	login string,
	email string,
) (string, error) {
	if s == nil || s.client == nil {
		return "", ErrNilFirebaseAuthClient
	}
	if githubID <= 0 {
		return "", ErrInvalidGitHubID
	}
	if login == "" {
		return "", ErrEmptyGitHubLogin
	}

	githubUID := strconv.FormatInt(githubID, 10)

	user, err := s.client.GetUserByProviderUID(
		ctx,
		githubProviderID,
		githubUID,
	)
	if err != nil {
		if !s.isUserNotFound(err) {
			return "", fmt.Errorf("get Firebase user by GitHub ID: %w", err)
		}

		user, err = s.createAndLinkGitHubUser(
			ctx,
			githubUID,
			login,
			email,
		)
		if err != nil {
			return "", err
		}

		return user.UID, nil
	}

	if user == nil || user.UID == "" {
		return "", ErrFirebaseUserMissingUID
	}

	if err := s.updateProfile(ctx, user.UID, login, email, nil); err != nil {
		return "", err
	}

	return user.UID, nil
}

func (s *Service) createAndLinkGitHubUser(
	ctx context.Context,
	githubUID string,
	login string,
	email string,
) (*auth.UserRecord, error) {
	toCreate := (&auth.UserToCreate{}).DisplayName(login)
	if email != "" {
		toCreate.Email(email)
	}

	user, err := s.client.CreateUser(ctx, toCreate)
	if err != nil {
		return nil, fmt.Errorf("create Firebase user: %w", err)
	}
	if user == nil || user.UID == "" {
		return nil, ErrFirebaseUserMissingUID
	}

	provider := &auth.UserProvider{
		UID:         githubUID,
		ProviderID:  githubProviderID,
		Email:       email,
		DisplayName: login,
	}

	if err := s.updateProfile(ctx, user.UID, login, email, provider); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) updateProfile(
	ctx context.Context,
	firebaseUID string,
	login string,
	email string,
	providerToLink *auth.UserProvider,
) error {
	toUpdate := (&auth.UserToUpdate{}).DisplayName(login)
	if email != "" {
		toUpdate.Email(email)
	}
	if providerToLink != nil {
		toUpdate.ProviderToLink(providerToLink)
	}

	if _, err := s.client.UpdateUser(ctx, firebaseUID, toUpdate); err != nil {
		return fmt.Errorf("update Firebase user: %w", err)
	}

	return nil
}
