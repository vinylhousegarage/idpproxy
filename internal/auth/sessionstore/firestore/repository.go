package firesessionstore

import (
	"cloud.google.com/go/firestore"
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
