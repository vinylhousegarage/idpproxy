package store

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
)

func (r *Repo) MarkUsed(ctx context.Context, refreshID string) error {
	refreshID = strings.TrimSpace(refreshID)

	if refreshID == "" {
		return fmt.Errorf("%w: empty", ErrInvalidID)
	}
	if err := validateRefreshID(refreshID); err != nil {
		return err
	}

	doc := r.docRT(refreshID)

	return r.fs.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		now := r.now()

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
		}, firestore.LastUpdateTime(snap.UpdateTime))
	})
}
