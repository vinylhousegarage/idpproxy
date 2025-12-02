package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUsecase_Invalidate(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		uc        *Usecase
		repo      *fakeRepository
		sessionID string
		wantErr   error
	}{
		{
			name:      "nil_usecase",
			uc:        nil,
			repo:      nil,
			sessionID: "session-123",
			wantErr:   ErrInvalidUsecaseConfig,
		},
		{
			name: "nil_repo",
			uc: &Usecase{
				Repo: nil,
				Now:  time.Now,
			},
			sessionID: "session-123",
			wantErr:   ErrInvalidUsecaseConfig,
		},
		{
			name: "nil_now",
			uc: &Usecase{
				Repo: newFakeRepository(),
				Now:  nil,
			},
			sessionID: "session-123",
			wantErr:   ErrInvalidUsecaseConfig,
		},
		{
			name: "empty_session_id",
			uc: &Usecase{
				Repo: newFakeRepository(),
				Now:  time.Now,
			},
			sessionID: "",
			wantErr:   ErrEmptySessionID,
		},
		{
			name: "find_error",
			uc: &Usecase{
				Repo: func() *fakeRepository {
					r := newFakeRepository()
					r.findErr = errors.New("find error")
					return r
				}(),
				Now: time.Now,
			},
			sessionID: "session-123",
			wantErr:   errors.New("find error"),
		},
		{
			name: "expired_session",
			uc: &Usecase{
				Repo: newFakeRepositoryWithSession(&Session{
					SessionID: "session-123",
					Status:    "active",
					ExpiresAt: now.Add(-time.Minute),
				}),
				Now: func() time.Time { return now },
			},
			sessionID: "session-123",
			wantErr:   ErrExpiredSession,
		},
		{
			name: "inactive_session",
			uc: &Usecase{
				Repo: newFakeRepositoryWithSession(&Session{
					SessionID: "session-123",
					Status:    "inactive",
					ExpiresAt: now.Add(time.Hour),
				}),
				Now: func() time.Time { return now },
			},
			sessionID: "session-123",
			wantErr:   ErrInactiveSession,
		},
		{
			name: "update_error",
			uc: &Usecase{
				Repo: func() *fakeRepository {
					r := newFakeRepositoryWithSession(&Session{
						SessionID: "session-123",
						Status:    "active",
						ExpiresAt: now.Add(time.Hour),
					})
					r.updateErr = errors.New("update error")
					return r
				}(),
				Now: func() time.Time { return now },
			},
			sessionID: "session-123",
			wantErr:   errors.New("update error"),
		},
		{
			name: "success",
			uc: &Usecase{
				Repo: newFakeRepositoryWithSession(&Session{
					SessionID: "session-123",
					UserID:    "user-123",
					Status:    "active",
					ExpiresAt: now.Add(time.Hour),
				}),
				Now: func() time.Time { return now },
			},
			sessionID: "session-123",
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			got, err := tt.uc.Invalidate(ctx, tt.sessionID)

			if tt.wantErr != nil {
				require.Nil(t, got)
				require.EqualError(t, err, tt.wantErr.Error())
				return
			}

			// success case
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, "inactive", got.Status)
			require.True(t, got.UpdatedAt.Equal(now))
		})
	}
}
