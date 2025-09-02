package store

import (
	"context"
	"time"
)

type GitHubTokenRepo interface {
	Upsert(ctx context.Context, rec *GitHubTokenRecord) error
	GetByFirebaseUID(ctx context.Context, uid string) (*GitHubTokenRecord, error)
	TouchLastUsed(ctx context.Context, uid string, t time.Time) error
	DeleteByFirebaseUID(ctx context.Context, uid string) error
}
