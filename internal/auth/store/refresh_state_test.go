package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIsActive(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_800_000_000, 0)

	tests := []struct {
		name string
		rec  *RefreshTokenRecord
		want bool
	}{
		{
			name: "nil record",
			rec:  nil,
			want: false,
		},
		{
			name: "active",
			rec: &RefreshTokenRecord{
				ExpiresAt: now.Add(1 * time.Hour),
			},
			want: true,
		},
		{
			name: "revoked",
			rec: &RefreshTokenRecord{
				RevokedAt: now.Add(-1 * time.Hour),
				ExpiresAt: now.Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "replaced",
			rec: &RefreshTokenRecord{
				ReplacedBy: "new-token",
				ExpiresAt:  now.Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "expired",
			rec: &RefreshTokenRecord{
				ExpiresAt: now.Add(-1 * time.Minute),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isActive(tt.rec, now)
			require.Equal(t, tt.want, got)
		})
	}
}
