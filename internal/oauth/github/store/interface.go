package store

import "context"

type GitHubTokenRepo interface {
	Upsert(ctx context.Context, rec *GitHubTokenRecord) error
	GetByFirebaseUID(ctx context.Context, uid string) (*GitHubTokenRecord, error)
	TouchLastUsed(ctx context.Context, uid string) error
	DeleteByFirebaseUID(ctx context.Context, uid string) error
}
