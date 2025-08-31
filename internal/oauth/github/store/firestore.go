package store

import (
	"time"

	"cloud.google.com/go/firestore"
)

const (
	collectionGitHubTokens = "github_tokens"
)

type FirestoreGitHubTokenRepo struct {
	col *firestore.CollectionRef
	now func() time.Time
	enc TokenEncryptor
}

func NewFirestoreGitHubTokenRepo(client *firestore.Client, enc TokenEncryptor) *FirestoreGitHubTokenRepo {
	return &FirestoreGitHubTokenRepo{
		col: client.Collection(collectionGitHubTokens),
		now: time.Now,
		enc: enc,
	}
}
