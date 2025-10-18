package store

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *Repo) Get(ctx context.Context, userID string) (*AccessGenerationRecord, error) {
	if userID == "" {
		return nil, ErrInvalidID
	}

	snap, err := r.docAG(userID).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var rec AccessGenerationRecord
	if err := snap.DataTo(&rec); err != nil {
		return nil, err
	}
	if rec.UserID == "" {
		rec.UserID = userID
	}

	return &rec, nil
}
