package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	colRefreshTokens     = "refresh_tokens"
	colAccessGenerations = "access_generations"
)

type RefreshRepo interface {
	Create(ctx context.Context, rec *RefreshTokenRecord) error
	GetByID(ctx context.Context, refreshID string) (*RefreshTokenRecord, error)
	MarkUsed(ctx context.Context, refreshID string, t time.Time) error
	Revoke(ctx context.Context, refreshID, reason string, t time.Time) error
	Replace(ctx context.Context, oldID string, newRec *RefreshTokenRecord, t time.Time) error
	RevokeFamily(ctx context.Context, familyID, reason string, t time.Time) (int, error)
	DeleteExpired(ctx context.Context, until time.Time) (int, error)
}

type AccessGenRepo interface {
	Get(ctx context.Context, userID string) (*AccessGenerationRecord, error)
	Set(ctx context.Context, rec *AccessGenerationRecord) error
	Bump(ctx context.Context, userID string, t time.Time) (newGen int, err error)
}

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

func (r *Repo) docAG(userID string) *firestore.DocumentRef {
	return r.fs.Collection(colAccessGenerations).Doc(userID)
}
