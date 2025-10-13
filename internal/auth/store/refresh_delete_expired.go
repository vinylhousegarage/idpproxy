package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (r *Repo) deleteByQuery(ctx context.Context, q firestore.Query) (int, error) {
	const batchSize = 500

	iter := q.Documents(ctx)
	defer iter.Stop()

	deleted := 0
	batch := r.fs.Batch()
	pending := 0

	commit := func() error {
		if pending == 0 {
			return nil
		}
		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
		deleted += pending

		batch = r.fs.Batch()
		pending = 0

		return nil
	}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return deleted, err
		}

		batch.Delete(doc.Ref)
		pending++
		if pending >= batchSize {
			if err := commit(); err != nil {
				return deleted, err
			}
		}
	}

	if err := commit(); err != nil {
		return deleted, err
	}

	return deleted, nil
}

func (r *Repo) DeleteExpired(ctx context.Context, until time.Time) (int, error) {
	if until.IsZero() {
		return 0, ErrInvalidUntil
	}

	col := r.fs.Collection(colRefreshTokens)
	total := 0

	if n, err := r.deleteByQuery(ctx,
		col.Where("expires_at", "<=", until),
	); err != nil {
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
