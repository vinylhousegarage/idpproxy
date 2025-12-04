package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUsecase_Touch(t *testing.T) {
	t.Parallel()

	t.Run("nil_usecase", func(t *testing.T) {
		t.Parallel()

		var uc *Usecase
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInvalidUsecaseConfig)
	})

	t.Run("nil_repo", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{
			Repo: nil,
			Now:  time.Now,
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInvalidUsecaseConfig)
	})

	t.Run("nil_now", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{
			Repo: newFakeRepository(),
			Now:  nil,
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInvalidUsecaseConfig)
	})

	t.Run("empty_session_id", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{
			Repo: newFakeRepository(),
			Now:  time.Now,
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrEmptySessionID)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)

		repo := newFakeRepository()
		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return now
			},
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "unknown-session")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrNotFound)
		require.Equal(t, "unknown-session", repo.lastFindID)
	})

	t.Run("expired_session", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)
		expired := now.Add(-time.Minute)

		s := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
			ExpiresAt: expired,
		}

		repo := newFakeRepositoryWithSession(s)
		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return now
			},
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrExpiredSession)
		require.Len(t, repo.updated, 0, "expired session must not be updated")
	})

	t.Run("inactive_session", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)
		future := now.Add(time.Minute)

		s := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "inactive",
			ExpiresAt: future,
		}

		repo := newFakeRepositoryWithSession(s)
		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return now
			},
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInactiveSession)
		require.Len(t, repo.updated, 0, "inactive session must not be updated")
	})

	t.Run("update_error", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)
		future := now.Add(time.Minute)

		s := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
			ExpiresAt: future,
		}

		repo := newFakeRepositoryWithSession(s)
		repo.updateErr = errors.New("update failed")

		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return now
			},
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, repo.updateErr)
		require.Len(t, repo.updated, 0, "on update error, updated slice should not be appended")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)
		future := now.Add(30 * time.Minute)

		s := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
			ExpiresAt: future,
		}

		repo := newFakeRepositoryWithSession(s)
		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return now
			},
		}
		ctx := context.Background()

		got, err := uc.Touch(ctx, "session-123")

		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, "session-123", got.SessionID)
		require.Equal(t, "user-123", got.UserID)
		require.Equal(t, "active", got.Status)

		require.NotNil(t, got.LastUsed)
		require.Equal(t, now, *got.LastUsed)

		require.NotNil(t, got.UpdatedAt)
		require.Equal(t, now, *got.UpdatedAt)

		require.Equal(t, future, got.ExpiresAt)

		require.Len(t, repo.updated, 1)
		require.Equal(t, got, repo.updated[0])

		require.Equal(t, "session-123", repo.lastFindID)
	})
}
