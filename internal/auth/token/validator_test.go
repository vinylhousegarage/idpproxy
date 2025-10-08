package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    string
		ttl       time.Duration
		purge     time.Duration
		wantError error
	}{
		{"ok", "u1", time.Hour, 2 * time.Hour, nil},
		{"empty user", "", time.Hour, 2 * time.Hour, ErrEmptyUserID},
		{"ttl zero", "u1", 0, time.Hour, ErrInvalidTTL},
		{"ttl negative", "u1", -time.Second, time.Hour, ErrInvalidTTL},
		{"purge < ttl", "u1", time.Hour, 30 * time.Minute, ErrInvalidPurge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParams(tt.userID, tt.ttl, tt.purge)
			if tt.wantError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.wantError)
			}
		})
	}
}
