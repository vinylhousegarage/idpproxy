package store

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type bumpFunc func(ctx context.Context, userID string, t time.Time) (int, error)

func bumpWithRetry(ctx context.Context, userID string, t time.Time, fn bumpFunc) (int, error) {
	const maxAttempt = 5
	backoff := 50 * time.Millisecond

	var lastErr error
	for attempt := 1; attempt <= maxAttempt; attempt++ {
		gen, err := fn(ctx, userID, t)
		if err == nil {
			return gen, nil
		}
		switch status.Code(err) {
		case codes.Aborted, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Unavailable:
			lastErr = err
		default:
			return 0, err
		}
		select {
		case <-time.After(backoff):
			backoff *= 2
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	return 0, lastErr
}

func (r *Repo) BumpWithRetry(ctx context.Context, userID string, t time.Time) (int, error) {
	return bumpWithRetry(ctx, userID, t, r.Bump)
}
