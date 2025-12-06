package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakePurgeRepository struct {
	purgeBeforeCalls []time.Time
	purgeResult      int
	purgeErr         error
}

func (f *fakePurgeRepository) Create(_ context.Context, _ *Session) error {
	return nil
}

func (f *fakePurgeRepository) FindByID(_ context.Context, _ string) (*Session, error) {
	return nil, nil
}

func (f *fakePurgeRepository) Update(_ context.Context, _ *Session) error {
	return nil
}

func (f *fakePurgeRepository) PurgeExpired(_ context.Context, before time.Time) (int, error) {
	f.purgeBeforeCalls = append(f.purgeBeforeCalls, before)

	if f.purgeErr != nil {
		return 0, f.purgeErr
	}

	return f.purgeResult, nil
}

func TestUsecase_PurgeExpired(t *testing.T) {
	t.Parallel()

	t.Run("nil_usecase", func(t *testing.T) {
		t.Parallel()

		var uc *Usecase
		ctx := context.Background()

		gotCount, err := uc.PurgeExpired(ctx)

		require.Zero(t, gotCount)
		require.ErrorIs(t, err, ErrInvalidUsecaseConfig)
	})

	t.Run("nil_repo", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{
			Repo: nil,
			Now:  time.Now,
		}
		ctx := context.Background()

		gotCount, err := uc.PurgeExpired(ctx)

		require.Zero(t, gotCount)
		require.ErrorIs(t, err, ErrInvalidUsecaseConfig)
	})

	t.Run("nil_now", func(t *testing.T) {
		t.Parallel()

		uc := &Usecase{
			Repo: &fakePurgeRepository{},
			Now:  nil,
		}
		ctx := context.Background()

		gotCount, err := uc.PurgeExpired(ctx)

		require.Zero(t, gotCount)
		require.ErrorIs(t, err, ErrInvalidUsecaseConfig)
	})

	t.Run("repo_error", func(t *testing.T) {
		t.Parallel()

		repoErr := errors.New("purge failed")
		repo := &fakePurgeRepository{
			purgeErr: repoErr,
		}
		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)
			},
		}
		ctx := context.Background()

		gotCount, err := uc.PurgeExpired(ctx)

		require.Zero(t, gotCount)
		require.Error(t, err)
		require.ErrorIs(t, err, repoErr)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)
		repo := &fakePurgeRepository{
			purgeResult: 3,
		}
		uc := &Usecase{
			Repo: repo,
			Now: func() time.Time {
				return now
			},
		}
		ctx := context.Background()

		gotCount, err := uc.PurgeExpired(ctx)

		require.NoError(t, err)
		require.Equal(t, 3, gotCount)

		require.Len(t, repo.purgeBeforeCalls, 1)
		require.Equal(t, now, repo.purgeBeforeCalls[0])
	})
}
