package session

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeGetRepository struct {
	session *Session
	err     error
	lastID  string
}

func (f *fakeGetRepository) Create(_ context.Context, _ *Session) error {
	return nil
}

func (f *fakeGetRepository) FindByID(_ context.Context, sessionID string) (*Session, error) {
	f.lastID = sessionID
	if f.err != nil {
		return nil, f.err
	}
	return f.session, nil
}

func TestUsecase_Get(t *testing.T) {
	t.Parallel()

	t.Run("nil_usecase", func(t *testing.T) {
		t.Parallel()

		var uc *Usecase
		ctx := context.Background()

		got, err := uc.Get(ctx, "session-123")

		require.Nil(t, got)
		require.EqualError(t, err, "session: invalid usecase configuration")
	})

	t.Run("nil_repo", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{}
		ctx := context.Background()

		got, err := uc.Get(ctx, "session-123")

		require.Nil(t, got)
		require.EqualError(t, err, "session: invalid usecase configuration")
	})

	t.Run("empty_sessionID", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGetRepository{}
		uc := &Usecase{Repo: repo}
		ctx := context.Background()

		got, err := uc.Get(ctx, "")

		require.Nil(t, got)
		require.EqualError(t, err, "session: empty sessionID")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := &Session{
			SessionID: "session-123",
			UserID:    "user-123",
			Status:    "active",
		}
		repo := &fakeGetRepository{
			session: s,
		}
		uc := &Usecase{Repo: repo}
		ctx := context.Background()

		got, err := uc.Get(ctx, "session-123")

		require.NoError(t, err)
		require.Same(t, s, got)
		require.Equal(t, "session-123", repo.lastID)
	})

	t.Run("repository_error", func(t *testing.T) {
		t.Parallel()

		repoErr := errors.New("find error")
		repo := &fakeGetRepository{
			err: repoErr,
		}
		uc := &Usecase{Repo: repo}
		ctx := context.Background()

		got, err := uc.Get(ctx, "session-123")

		require.Nil(t, got)
		require.ErrorIs(t, err, repoErr)
		require.Equal(t, "session-123", repo.lastID)
	})
}
