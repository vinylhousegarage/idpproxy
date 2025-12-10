package firesessionstore

import (
	"context"

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

// var _ session.Repository = (*Repository)(nil)

func (r *Repository) Create(ctx context.Context, s *session.Session) error {
	_, err := r.collection.Doc(s.SessionID).Set(ctx, s)
	return err
}
