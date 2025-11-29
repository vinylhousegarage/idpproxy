// usecase_validate_test.go
package session

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUsecase_Validate(t *testing.T) {
	t.Parallel()

	fixedNow := time.Date(2025, 11, 29, 12, 0, 0, 0, time.UTC)

	t.Run("nil_usecase", func(t *testing.T) {
		t.Parallel()

		var uc *Usecase
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.Nil(t, got)
		require.EqualError(t, err, "session: invalid usecase configuration")
	})

	t.Run("nil_repo", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{
			Now: func() time.Time { return fixedNow },
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.Nil(t, got)
		require.EqualError(t, err, "session: invalid usecase configuration")
	})

	t.Run("nil_now", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepository()
		uc := &Usecase{
			Repo: repo,
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.Nil(t, got)
		require.EqualError(t, err, "session: invalid usecase configuration")
	})

	t.Run("empty_sessionID", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepository()
		uc := &Usecase{
			Repo: repo,
			Now:  func() time.Time { return fixedNow },
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "")

		require.Nil(t, got)
		require.EqualError(t, err, "session: empty sessionID")
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepository()
		repo.findErr = ErrNotFound

		uc := &Usecase{
			Repo: repo,
			Now:  func() time.Time { return fixedNow },
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("expired_session", func(t *testing.T) {
		t.Parallel()

		expired := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
			ExpiresAt: fixedNow.Add(-1 * time.Minute),
		}

		repo := newFakeRepository()
		repo.findMap["session-123"] = expired

		uc := &Usecase{
			Repo: repo,
			Now:  func() time.Time { return fixedNow },
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrExpiredSession)
	})

	t.Run("inactive_session", func(t *testing.T) {
		t.Parallel()

		inactive := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "revoked",
			ExpiresAt: fixedNow.Add(10 * time.Minute),
		}

		repo := newFakeRepository()
		repo.findMap["session-123"] = inactive

		uc := &Usecase{
			Repo: repo,
			Now:  func() time.Time { return fixedNow },
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInactiveSession)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		active := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
			ExpiresAt: fixedNow.Add(10 * time.Minute),
		}

		repo := newFakeRepository()
		repo.findMap["session-123"] = active

		uc := &Usecase{
			Repo: repo,
			Now:  func() time.Time { return fixedNow },
		}
		ctx := context.Background()

		got, err := uc.Validate(ctx, "session-123")

		require.NoError(t, err)
		require.Same(t, active, got)
	})
}
