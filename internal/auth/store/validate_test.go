package store

import (
	"errors"
	"strings"
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

	cases := []struct {
		name string
		in   string
		ok   bool
	}{
		{"ok-github", "github:" + okUUID, true},
		{"ok-google", "google:" + okUUID, true},
		{"ok-trim-and-case", "  GitHub :  " + okUUID + "  ", true},

		{"ng-empty", "   ", false},
		{"ng-no-colon", "github" + okUUID, false},
		{"ng-unknown-provider", "x:" + okUUID, false},
		{"ng-bad-uuid", "github:not-a-uuid", false},
		{"ng-missing-uuid", "github:", false},
		{"ng-missing-provider", ":" + okUUID, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validateUserID(tc.in)
			if tc.ok {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorIs(t, err, ErrInvalidUserID)
				if strings.Contains(tc.name, "no-colon") {
					require.Contains(t, err.Error(), "want '<provider>:<uuid>'")
				}
			}
		})
	}
}
