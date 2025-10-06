package store

import (
	"errors"
	"testing"

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
