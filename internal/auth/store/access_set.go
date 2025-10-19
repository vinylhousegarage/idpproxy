package store

import (
	"context"

	"cloud.google.com/go/firestore"
)

func (r *Repo) Set(ctx context.Context, rec *AccessGenerationRecord) error {
	if rec == nil {
		return ErrInvalidArgument
	}
	if rec.UserID == "" {
		return ErrInvalidID
	}

	rec.UpdatedAt = r.now().UTC()
	_, err := r.docAG(rec.UserID).Set(ctx, rec, firestore.MergeAll)

	return err
}
