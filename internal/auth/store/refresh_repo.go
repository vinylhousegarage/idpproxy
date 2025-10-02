package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateForCreate(rec *RefreshTokenRecord) error {
	if rec == nil {
		return fmt.Errorf("nil RefreshTokenRecord")
	}
	if rec.RefreshID == "" {
		return fmt.Errorf("missing RefreshID")
	}
	if strings.Contains(rec.RefreshID, "/") {
		return fmt.Errorf("invalid RefreshID: must not contain '/'")
	}
	if rec.UserID == "" {
		return fmt.Errorf("missing UserID")
	}
	if rec.DigestB64 == "" {
		return fmt.Errorf("missing DigestB64")
	}
	if rec.ReplacedBy != "" {
		return fmt.Errorf("replaced_by must be empty on create")
	}
	if !rec.RevokedAt.IsZero() {
		return fmt.Errorf("revoked_at must be zero on create")
	}

	return nil
}

func prepareForCreate(rec *RefreshTokenRecord, now time.Time) error {
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = now
	}
	if !rec.LastUsedAt.IsZero() && rec.LastUsedAt.Before(rec.CreatedAt) {
		return fmt.Errorf("last_used_at must be >= created_at")
	}
	if !rec.ExpiresAt.IsZero() && rec.ExpiresAt.Before(rec.CreatedAt) {
		return fmt.Errorf("expires_at must be >= created_at")
	}
	if !rec.DeleteAt.IsZero() {
		anchor := rec.ExpiresAt
		if anchor.IsZero() {
			anchor = rec.CreatedAt
		}
		if rec.DeleteAt.Before(anchor) {
			return fmt.Errorf("delete_at must be >= max(expires_at, created_at)")
		}
	}

	return nil
}

func (r *Repo) Create(ctx context.Context, rec *RefreshTokenRecord) error {
	if err := validateForCreate(rec); err != nil {
		return fmt.Errorf("create: validate: %w", err)
	}
	if err := prepareForCreate(rec, r.now()); err != nil {
		return fmt.Errorf("create: prepare: %w", err)
	}
	_, err := r.docRT(rec.RefreshID).Create(ctx, rec)
	if status.Code(err) == codes.AlreadyExists {
		return ErrConflict
	}

	return err
}
