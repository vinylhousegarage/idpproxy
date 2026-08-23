package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionGitHubTokens = "github_tokens"
)

type FirestoreGitHubTokenRepo struct {
	client *firestore.Client
	col    *firestore.CollectionRef
	now    func() time.Time
	enc    TokenEncryptor
}

type firestoreGitHubToken struct {
	GitHubID    string     `firestore:"github_id"`
	Provider    string     `firestore:"provider"`
	FirebaseUID string     `firestore:"firebase_uid"`
	Login       string     `firestore:"login"`
	Scopes      []string   `firestore:"scopes"`
	TokenType   string     `firestore:"token_type"`
	AccessToken Ciphertext `firestore:"access_token"`
	ExpiresAt   time.Time  `firestore:"expires_at"`
	LastUsedAt  time.Time  `firestore:"last_used_at"`
	CreatedAt   time.Time  `firestore:"created_at"`
	UpdatedAt   time.Time  `firestore:"updated_at"`
	DeleteAt    time.Time  `firestore:"delete_at"`
}

func NewFirestoreGitHubTokenRepo(
	client *firestore.Client,
	enc TokenEncryptor,
) *FirestoreGitHubTokenRepo {
	return &FirestoreGitHubTokenRepo{
		client: client,
		col:    client.Collection(collectionGitHubTokens),
		now:    time.Now,
		enc:    enc,
	}
}

func (r *FirestoreGitHubTokenRepo) Upsert(
	ctx context.Context,
	rec *GitHubTokenRecord,
) error {
	encryptedToken, err := r.enc.EncryptString(rec.AccessToken)
	if err != nil {
		return fmt.Errorf("encrypt GitHub access token: %w", err)
	}

	now := r.now()
	doc := r.col.Doc(rec.FirebaseUID)

	err = r.client.RunTransaction(ctx, func(
		ctx context.Context,
		tx *firestore.Transaction,
	) error {
		createdAt := now

		snapshot, err := tx.Get(doc)
		switch {
		case err == nil:
			var existing firestoreGitHubToken

			if err := snapshot.DataTo(&existing); err != nil {
				return fmt.Errorf("decode existing GitHub token: %w", err)
			}

			createdAt = existing.CreatedAt

		case status.Code(err) == codes.NotFound:

		default:
			return fmt.Errorf("get existing GitHub token: %w", err)
		}

		token := firestoreGitHubToken{
			GitHubID:    rec.GitHubID,
			Provider:    rec.Provider,
			FirebaseUID: rec.FirebaseUID,
			Login:       rec.Login,
			Scopes:      rec.Scopes,
			TokenType:   rec.TokenType,
			AccessToken: encryptedToken,
			ExpiresAt:   rec.ExpiresAt,
			LastUsedAt:  rec.LastUsedAt,
			CreatedAt:   createdAt,
			UpdatedAt:   now,
			DeleteAt:    rec.DeleteAt,
		}

		return tx.Set(doc, token)
	})
	if err != nil {
		return fmt.Errorf("upsert GitHub token: %w", err)
	}

	return nil
}
