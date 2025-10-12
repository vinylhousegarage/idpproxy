package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (r *Repo) RevokeFamily(ctx context.Context, familyID, reason string, t time.Time) (int, error) {
	if err := validateFamilyID(familyID); err != nil {
		return 0, fmt.Errorf("invalid familyID: %w", err)
	}
	if t.IsZero() {
		t = r.now()
	}

	iter := r.fs.Collection(colRefreshTokens).Where("family_id", "==", familyID).Documents(ctx)
	defer iter.Stop()

	bw := r.fs.BulkWriter(ctx)
	var jobs []*firestore.BulkWriterJob

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			bw.End()
			return 0, err
		}

		var rec RefreshTokenRecord
		if err := doc.DataTo(&rec); err != nil {
			bw.End()
			return 0, fmt.Errorf("decode: %w", err)
		}
		if !rec.RevokedAt.IsZero() {
			continue
	}
		j, err := bw.Update(doc.Ref, []firestore.Update{
			{Path: "revoked_at", Value: t},
			{Path: "revoke_reason", Value: reason},
		})
		if err != nil {
			bw.End()
			return 0, err
		}
		jobs = append(jobs, j)
	}

	bw.End()

	affected := 0
	var firstErr error
	for _, job := range jobs {
		if _, err := job.Results(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		affected++
	}

	return affected, firstErr
}
