package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const maxBatchWrites = 500

func (r *Repo) RevokeFamily(ctx context.Context, familyID, reason string, t time.Time) (int, error) {
	if err := validateFamilyID(familyID); err != nil {
		return 0, fmt.Errorf("invalid familyID: %w", err)
	}
	if t.IsZero() {
		t = r.now()
	}

	q := r.fs.Collection(colRefreshTokens).Where("family_id", "==", familyID)
	iter := q.Documents(ctx)
	defer iter.Stop()

	commit := func(b *firestore.WriteBatch) error {
		_, err := b.Commit(ctx)
		return err
	}

	batch := r.fs.Batch()
	pending := 0
	affected := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return affected, err
		}

		var rec RefreshTokenRecord
		if err := doc.DataTo(&rec); err != nil {
			return affected, fmt.Errorf("decode: %w", err)
		}

		if !rec.RevokedAt.IsZero() {
			continue
		}

		updates := []firestore.Update{
			{Path: "revoked_at", Value: t},
			{Path: "revoke_reason", Value: reason},
		}

		batch.Update(doc.Ref, updates)
		pending++
		affected++

		if pending >= maxBatchWrites-1 {
			if err := commit(batch); err != nil {
				return affected, err
			}
			batch = r.fs.Batch()
			pending = 0
		}
	}

	if pending > 0 {
		if err := commit(batch); err != nil {
			return affected, err
		}
	}

	return affected, nil
}
