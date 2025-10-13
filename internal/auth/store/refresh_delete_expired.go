package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (r *Repo) deleteByQuery(ctx context.Context, q firestore.Query) (int, error) {
	iter := q.Documents(ctx)
	defer iter.Stop()

	bw := r.fs.BulkWriter(ctx)
	defer bw.End()

	deleted := 0
	var jobs []*firestore.BulkWriterJob

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return deleted, err
		}

		job, err := bw.Delete(doc.Ref)
		if err != nil {
			return deleted, err
		}
		jobs = append(jobs, job)
	}

	bw.End()

	for _, j := range jobs {
		if _, err := j.Results(); err != nil {
			continue
		}
		deleted++
	}

	return deleted, nil
}

func (r *Repo) DeleteExpired(ctx context.Context, until time.Time) (int, error) {
	if until.IsZero() {
		return 0, ErrInvalidUntil
	}

	col := r.fs.Collection(colRefreshTokens)
	total := 0

	if n, err := r.deleteByQuery(ctx, col.Where("expires_at", "<=", until)); err != nil {
		return total, err
	} else {
		total += n
	}

	epoch := time.Unix(0, 0).UTC()
	if n, err := r.deleteByQuery(ctx,
		col.Where("revoked_at", ">", epoch).Where("revoked_at", "<=", until),
	); err != nil {
		return total, err
	} else {
		total += n
	}

	return total, nil
}
