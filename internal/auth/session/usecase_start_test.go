package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUsecase_Start(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		repoErr      error
		idGenErr     error
		wantErr      bool
		checkSession func(t *testing.T, repo *fakeRepository, got *Session, now time.Time, ttl time.Duration)
	}{
		{
			name:   "success",
			userID: "user-123",
			checkSession: func(t *testing.T, repo *fakeRepository, got *Session, now time.Time, ttl time.Duration) {
				require.Len(t, repo.created, 1)
				saved := repo.created[0]

				require.Equal(t, "session-123", saved.SessionID)
				require.Equal(t, "user-123", saved.UserID)
				require.Equal(t, "active", saved.Status)

				require.Equal(t, now, saved.CreatedAt)
				require.Equal(t, now.Add(ttl), saved.ExpiresAt)

				require.Nil(t, saved.UpdatedAt)
				require.Nil(t, saved.LastUsed)

				require.Equal(t, saved.SessionID, got.SessionID)
				require.Equal(t, saved.UserID, got.UserID)
			},
		},
		{
			name:    "empty userID",
			userID:  "",
			wantErr: true,
		},
		{
			name:     "IDGenerator error",
			userID:   "user-123",
			idGenErr: errors.New("idgen error"),
			wantErr:  true,
		},
		{
			name:    "Repo.Create error",
			userID:  "user-123",
			repoErr: errors.New("create failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			now := time.Date(2025, 11, 22, 10, 0, 0, 0, time.UTC)
			ttl := 30 * time.Minute

			repo := newFakeRepository()
			repo.createErr = tt.repoErr

			uc := &Usecase{
				Repo: repo,
				Now: func() time.Time {
					return now
				},
				TTL: ttl,
				IDGenerator: func() (string, error) {
					if tt.idGenErr != nil {
						return "", tt.idGenErr
					}
					return "session-123", nil
				},
			}

			got, err := uc.Start(context.Background(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			if tt.checkSession != nil {
				tt.checkSession(t, repo, got, now, ttl)
			}
		})
	}
}
