package firesessionstore

import (
	"cloud.google.com/go/firestore"
	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

type Repository struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
}

func NewRepository(client *firestore.Client, collectionName string) *Repository {
	return &Repository{
		client:     client,
		collection: client.Collection(collectionName),
	}
}

var _ session.Repository = (*Repository)(nil)
