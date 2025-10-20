package store

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *Repo) Bump(ctx context.Context, userID string, t time.Time) (int, error) {
	if userID == "" {
		return 0, ErrInvalidID
	}
	doc := r.docAG(userID)

	var newGen int

	err := r.fs.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(doc)
		switch status.Code(err) {
		case codes.NotFound:
			newGen = 1
			return tx.Set(doc, map[string]any{
				"user_id":    userID,
				"gen":        newGen,
				"updated_at": t.UTC(),
			})
		case codes.OK:

			var curGen int
			if v, err := snap.DataAt("gen"); err == nil {
				switch x := v.(type) {
				case int64:
					curGen = int(x)
				case float64:
					curGen = int(x)
				default:
					curGen = 0
				}
			}
			newGen = curGen + 1
			return tx.Update(doc, []firestore.Update{
				{Path: "gen", Value: firestore.Increment(1)},
				{Path: "updated_at", Value: t.UTC()},
			})
		default:
			return err
		}
	})
	if err != nil {
		return 0, err
	}

	return newGen, nil
}
