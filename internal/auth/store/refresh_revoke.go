package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

func (r *Repo) Revoke(ctx context.Context, id, reason string, t time.Time) error {
	return r.fs.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := r.docRT(id)
		snap, err := tx.Get(doc)
		if err != nil {
			return mapNotFound(err)
		}

		var rec RefreshTokenRecord
		if err := snap.DataTo(&rec); err != nil {
			return err
		}
		if !isActive(&rec, t) {
			return ErrAlreadyRevoked
		}

		updates := []firestore.Update{
			{Path: "revoked_at", Value: t},
			{Path: "revoke_reason", Value: reason},
		}

		return tx.Update(doc, updates, firestore.LastUpdateTime(snap.UpdateTime))
	})
}
