package fireauthcodestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
	authcodestore "github.com/vinylhousegarage/idpproxy/internal/authcode/store"
)

type Repository struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
	now        func() time.Time
}

func NewRepository(
	client *firestore.Client,
	collectionName string,
) *Repository {
	return &Repository{
		client:     client,
		collection: client.Collection(collectionName),
		now:        time.Now,
	}
}

var _ authcodestore.Store = (*Repository)(nil)

func (r *Repository) Save(
	ctx context.Context,
	proxyCode authcode.ProxyCode,
) error {
	_, err := r.collection.Doc(proxyCode.Code).Set(ctx, proxyCode)

	return err
}

func (r *Repository) Consume(
	ctx context.Context,
	proxyCodeValue string,
	clientID string,
) (authcode.ProxyCode, error) {
	var proxyCode authcode.ProxyCode
	var expired bool

	ref := r.collection.Doc(proxyCodeValue)

	err := r.client.RunTransaction(
		ctx,
		func(
			ctx context.Context,
			tx *firestore.Transaction,
		) error {
			expired = false

			doc, err := tx.Get(ref)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					return authcodestore.ErrNotFound
				}

				return err
			}

			if err := doc.DataTo(&proxyCode); err != nil {
				return err
			}

			if proxyCode.ClientID != clientID {
				return authcodestore.ErrClientMismatch
			}

			if !r.now().Before(proxyCode.ExpiresAt) {
				expired = true

				return tx.Delete(ref)
			}

			return tx.Delete(ref)
		},
	)
	if err != nil {
		return authcode.ProxyCode{}, err
	}

	if expired {
		return authcode.ProxyCode{}, authcodestore.ErrExpired
	}

	return proxyCode, nil
}
