package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/model"
)

type GoogleTokenStore interface {
	SaveGoogleToken(ctx context.Context, token *model.GoogleToken) (string, error)
	FindGoogleTokenByID(ctx context.Context, id string) (*model.GoogleToken, error)
	DeleteGoogleTokenByID(ctx context.Context, id string) error
}

type GoogleTokenRepository struct {
	Client *firestore.Client
	Logger *zap.Logger
}

func NewGoogleTokenRepository(client *firestore.Client, logger *zap.Logger) *GoogleTokenRepository {
	return &GoogleTokenRepository{Client: client, Logger: logger}
}

func (r *GoogleTokenRepository) collection() *firestore.CollectionRef {
	return r.Client.Collection(config.CollectionGoogleTokens)
}

func (r *GoogleTokenRepository) SaveGoogleToken(ctx context.Context, token *model.GoogleToken) (string, error) {
	docRef, _, err := r.collection().Add(ctx, token)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to save Google token to Firestore: %w", err)
		r.Logger.Error("failed to save Google token to Firestore", zap.Error(err))
		return "", wrappedErr
	}
	return docRef.ID, nil
}

func (r *GoogleTokenRepository) FindGoogleTokenByID(ctx context.Context, id string) (*model.GoogleToken, error) {
	docSnap, err := r.collection().Doc(id).Get(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to find Google token by ID (%s): %w", id, err)
		r.Logger.Error("failed to find Google token by ID", zap.String("id", id), zap.Error(err))
		return nil, wrappedErr
	}

	var token model.GoogleToken
	if err := docSnap.DataTo(&token); err != nil {
		wrappedErr := fmt.Errorf("failed to decode Firestore document to GoogleToken: %w", err)
		r.Logger.Error("failed to decode Firestore document to GoogleToken", zap.Error(err))
		return nil, wrappedErr
	}
	return &token, nil
}

func (r *GoogleTokenRepository) DeleteGoogleTokenByID(ctx context.Context, id string) error {
	_, err := r.collection().Doc(id).Delete(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete Google token by ID (%s): %w", id, err)
		r.Logger.Error("failed to delete Google token by ID", zap.String("id", id), zap.Error(err))
		return wrappedErr
	}
	return nil
}
