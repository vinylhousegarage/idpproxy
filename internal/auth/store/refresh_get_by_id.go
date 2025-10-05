package store

import (
	"context"
	"fmt"
	"strings"
)

func validateRefreshID(id string) error {
	if id == "" {
		return fmt.Errorf("%w: empty", ErrInvalidID)
	}
	if strings.Contains(id, "/") {
		return fmt.Errorf("%w: must not contain '/'", ErrInvalidID)
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, refreshID string) (*RefreshTokenRecord, error) {
	if err := validateRefreshID(refreshID); err != nil {
		return nil, err
	}

	snap, err := r.docRT(refreshID).Get(ctx)
	if err != nil {
		return nil, mapNotFound(err)
	}

	var rec RefreshTokenRecord
	if err := snap.DataTo(&rec); err != nil {
		return nil, err
	}

	rec.RefreshID = snap.Ref.ID
	return &rec, nil
}
