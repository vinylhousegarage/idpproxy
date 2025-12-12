package firesessionstore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (r *Repository) FindByID(ctx context.Context, sessionID string) (*session.Session, error) {
	doc, err := r.collection.Doc(sessionID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, session.ErrNotFound
		}
		return nil, err
	}

	var s session.Session
	if err := doc.DataTo(&s); err != nil {
		return nil, err
	}

	return &s, nil
}
