package store

import (
	"time"

	"cloud.google.com/go/firestore"
)

const (
	colRefreshTokens     = "refresh_tokens"
	colAccessGenerations = "access_generations"
)

type Repo struct {
	fs  *firestore.Client
	now func() time.Time
}

func NewRepo(fs *firestore.Client) *Repo {
	return &Repo{fs: fs, now: time.Now}
}

func (r *Repo) docRT(id string) *firestore.DocumentRef {
	return r.fs.Collection(colRefreshTokens).Doc(id)
}
