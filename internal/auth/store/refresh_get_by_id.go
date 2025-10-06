package store

import (
	"context"
)

func (r *Repo) GetByID(ctx context.Context, refreshID string) (*RefreshTokenRecord, error) {
	if err := validateRefreshID(refreshID); err != nil {
		return nil, err
	}

	snap, err := r.docRT(refreshID).Get(ctx)
	if err != nil {
		return nil, mapNotFound(err)
	}

	var rec RefreshTokenRecord
	if err := snap.DataTo(&rec); err != nil {
		return nil, err
	}

	rec.RefreshID = snap.Ref.ID
	return &rec, nil
}
