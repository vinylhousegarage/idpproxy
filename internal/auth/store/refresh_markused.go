package store

import (
	"context"
	"strings"

	"cloud.google.com/go/firestore"
)

func validateRefreshID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" || strings.Contains(id, "/") {
		return ErrInvalidID
	}

	return nil
}

func (r *Repo) MarkUsed(ctx context.Context, refreshID string) error {
	if err := validateRefreshID(refreshID); err != nil {
		return err
	}

	doc := r.docRT(refreshID)
	now := r.now()

	return r.fs.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(doc)
		if err != nil {
			return mapNotFound(err)
		}

		var rec RefreshTokenRecord
		if err := snap.DataTo(&rec); err != nil {
			return err
		}

		if !rec.RevokedAt.IsZero() {
			return ErrRevoked
		}
		if !rec.DeleteAt.IsZero() && !now.Before(rec.DeleteAt) {
			return ErrDeleted
		}
		if !rec.ExpiresAt.IsZero() && !now.Before(rec.ExpiresAt) {
			return ErrExpired
		}

		return tx.Update(doc, []firestore.Update{
			{Path: "last_used_at", Value: now},
		})
	})
}
