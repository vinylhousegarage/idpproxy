package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateReplaceArgs(oldID string, newRec *RefreshTokenRecord) error {
	if newRec == nil {
		return ErrInvalid
	}
	if err := validateRefreshID(strings.TrimSpace(oldID)); err != nil {
		return fmt.Errorf("oldID: %w", err)
	}
	if err := validateRefreshID(strings.TrimSpace(newRec.RefreshID)); err != nil {
		return fmt.Errorf("new RefreshID: %w", err)
	}
	if oldID == newRec.RefreshID {
		return ErrConflict
	}
	if err := validateUserID(strings.TrimSpace(newRec.UserID)); err != nil {
		return fmt.Errorf("new UserID: %w", err)
	}

	if newRec.DigestB64 == "" {
		return ErrInvalid
	}

	if !newRec.RevokedAt.IsZero() {
		return ErrInvalid
	}
	if newRec.ReplacedBy != "" {
		return ErrInvalid
	}

	return nil
}

func checkReplaceAllowed(old *RefreshTokenRecord, newRec *RefreshTokenRecord, t time.Time) error {
	if !isActive(old, t) {
		return ErrConflict
	}
	if old.UserID != newRec.UserID {
		return ErrConflict
	}
	return nil
}

func prepareNewFromOld(dst *RefreshTokenRecord, src *RefreshTokenRecord, t time.Time) {
	dst.FamilyID = src.FamilyID
	dst.CreatedAt = t
}

func txUpdateWithPrecond(tx *firestore.Transaction, ref *firestore.DocumentRef, snap *firestore.DocumentSnapshot, ups []firestore.Update) error {
	return tx.Update(ref, ups, firestore.LastUpdateTime(snap.UpdateTime))
}

func (r *Repo) Replace(ctx context.Context, oldID string, newRec *RefreshTokenRecord, t time.Time) error {
	if err := validateReplaceArgs(oldID, newRec); err != nil {
		return err
	}

	err := r.fs.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		oldRef := r.docRT(oldID)
		oldSnap, err := tx.Get(oldRef)
		if err != nil {
			return mapNotFound(err)
		}

		var old RefreshTokenRecord
		if err := oldSnap.DataTo(&old); err != nil {
			return err
		}

		if err := checkReplaceAllowed(&old, newRec, t); err != nil {
			return err
		}
		prepareNewFromOld(newRec, &old, t)

		newRef := r.docRT(newRec.RefreshID)

		if _, err := tx.Get(newRef); err == nil {
			return ErrConflict
		} else if status.Code(err) != codes.NotFound {
			return err
		}

		if err := txUpdateWithPrecond(tx, oldRef, oldSnap, []firestore.Update{
			{Path: "replaced_by", Value: newRec.RefreshID},
			{Path: "revoked_at", Value: t},
		}); err != nil {
			return mapConflict(err)
		}

		if err := tx.Create(newRef, newRec); err != nil {
			return mapConflict(err)
		}
		return nil
	})

	return mapConflict(err)
}
