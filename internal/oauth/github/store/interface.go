package store

import (
	"context"
	"time"
)

type TokenEncryptor interface {
	EncryptString(plain string) (string, error)
	DecryptString(token string) (string, error)
}

type GitHubTokenRepo interface {
	Upsert(ctx context.Context, rec *GitHubTokenRecord) error
	GetByFirebaseUID(ctx context.Context, uid string) (*GitHubTokenRecord, error)
	TouchLastUsed(ctx context.Context, uid string, t time.Time) error
	DeleteByFirebaseUID(ctx context.Context, uid string) error
}
