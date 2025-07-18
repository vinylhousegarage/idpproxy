package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/google/model"
)

type GoogleTokenRepository struct {
	Client *firestore.Client
	Logger *zap.Logger
}

func NewGoogleTokenRepository(client *firestore.Client, logger *zap.Logger) *GoogleTokenRepository {
	return &GoogleTokenRepository{Client: client, Logger: logger}
}

func (r *GoogleTokenRepository) SaveGoogleToken(ctx context.Context, token *model.GoogleToken) (string, error) {
	docRef, _, err := r.Client.Collection(config.CollectionGoogleTokens).Add(ctx, token)
	if err != nil {
		r.Logger.Error("failed to save Google token to Firestore", zap.Error(err))
		return "", err
	}
	return docRef.ID, nil
}
