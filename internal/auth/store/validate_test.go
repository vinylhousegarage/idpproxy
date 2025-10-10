package store

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateRefreshID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"ok-basic", "abc123", false},
		{"ok-trim", "  abc  ", false},
		{"err-empty", "   ", true},
		{"err-slash", "a/b", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateRefreshID(tt.id)
			if tt.wantErr {
				require.Error(t, err)
				require.True(t, errors.Is(err, ErrInvalidID), "err=%v", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserID(t *testing.T) {
	t.Parallel()

	okUUID := uuid.New().String()

	tests := []struct {
		name   string
		userID string
		wantOK bool
	}{
		{"github-ok", "github:" + okUUID, true},
		{"google-ok", "google:" + okUUID, true},
		{"missing", "", false},
		{"no-colon", "github" + okUUID, false},
		{"unknown-provider", "x:" + okUUID, false},
		{"bad-uuid", "github:not-a-uuid", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateUserID(tt.userID)
			if tt.wantOK && err != nil {
				t.Fatalf("want OK, got err: %v", err)
			}
			if !tt.wantOK && err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}
}
