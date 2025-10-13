package store

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRepo_DeleteExpired(t *testing.T) {
	requireEmulator(t)

	t.Run("until is zero -> ErrInvalidUntil", func(t *testing.T) {
		r := newTestRepo(t)
		_, err := r.DeleteExpired(context.Background(), time.Time{})
		if !errors.Is(err, ErrInvalidUntil) {
			t.Fatalf("want ErrInvalidUntil, got %v", err)
		}
	})

	t.Run("delete only expired and revoked (keep active)", func(t *testing.T) {
		fixed := time.Unix(1_800_000_000, 0).UTC()
		r := newTestRepoWithNow(t, fixed)

		purgeRefreshTokens(t, r)

		until := fixed

		expired := makeActiveRec("rt-expired-1", "github:11111111-1111-1111-1111-111111111111", fixed)
		expired.RefreshID = "rt-expired-1"
		expired.ExpiresAt = fixed.Add(-time.Hour)
		expired.DeleteAt = fixed.Add(24 * time.Hour)
		seedRefreshDoc(t, r, expired)

		revoked := makeActiveRec("rt-revoked-1", "github:22222222-2222-2222-2222-222222222222", fixed)
		revoked.RefreshID = "rt-revoked-1"
		revoked.RevokedAt = fixed.Add(-30 * time.Minute)
		revoked.DeleteAt = fixed.Add(24 * time.Hour)
		seedRefreshDoc(t, r, revoked)

		active := makeActiveRec("rt-active-1", "github:33333333-3333-3333-3333-333333333333", fixed)
		active.RefreshID = "rt-active-1"
		active.ExpiresAt = fixed.Add(24 * time.Hour)
		active.RevokedAt = time.Time{}
		active.DeleteAt = fixed.Add(30 * 24 * time.Hour)
		seedRefreshDoc(t, r, active)

		futureExp := makeActiveRec("rt-future-exp-1", "github:44444444-4444-4444-4444-444444444444", fixed)
		futureExp.RefreshID = "rt-future-exp-1"
		futureExp.ExpiresAt = fixed.Add(2 * time.Hour)
		futureExp.RevokedAt = time.Time{}
		futureExp.DeleteAt = fixed.Add(30 * 24 * time.Hour)
		seedRefreshDoc(t, r, futureExp)

		got, err := r.DeleteExpired(context.Background(), until)
		if err != nil {
			t.Fatalf("DeleteExpired error: %v", err)
		}
		if got != 2 {
			t.Fatalf("want deleted=2, got=%d", got)
		}

		ctx := context.Background()

		if _, err := r.docRT("rt-expired-1").Get(ctx); status.Code(err) != codes.NotFound {
			t.Fatalf("expired should be deleted; got err=%v", err)
		}
		if _, err := r.docRT("rt-revoked-1").Get(ctx); status.Code(err) != codes.NotFound {
			t.Fatalf("revoked should be deleted; got err=%v", err)
		}
		if _, err := r.docRT("rt-active-1").Get(ctx); err != nil {
			t.Fatalf("active should remain; got err=%v", err)
		}
		if _, err := r.docRT("rt-future-exp-1").Get(ctx); err != nil {
			t.Fatalf("future-exp should remain; got err=%v", err)
		}
	})

	t.Run("batch deletion over 500 docs (e.g., 503)", func(t *testing.T) {
		fixed := time.Unix(1_800_100_000, 0).UTC()
		r := newTestRepoWithNow(t, fixed)

		purgeRefreshTokens(t, r)

		until := fixed

		n := 503
		for i := 0; i < n; i++ {
			id := fmt.Sprintf("rt-batch-exp-%03d", i)
			rec := makeActiveRec(id, "github:55555555-5555-5555-5555-555555555555", fixed)
			rec.ExpiresAt = fixed.Add(-time.Minute)
			rec.DeleteAt = fixed.Add(24 * time.Hour)
			seedRefreshDoc(t, r, rec)
		}

		deleted, err := r.DeleteExpired(context.Background(), until)
		if err != nil {
			t.Fatalf("DeleteExpired error: %v", err)
		}
		if deleted != n {
			t.Fatalf("want deleted=%d, got=%d", n, deleted)
		}

		ctx := context.Background()
		for _, probe := range []string{"rt-batch-exp-000", "rt-batch-exp-250", "rt-batch-exp-502"} {
			if _, err := r.docRT(probe).Get(ctx); status.Code(err) != codes.NotFound {
				t.Fatalf("doc %q should be deleted; got err=%v", probe, err)
			}
		}
	})
}
