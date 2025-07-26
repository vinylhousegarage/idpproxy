package loginfirebase

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		wantErr   error
		wantToken string
	}{
		{
			name:      "valid request",
			body:      `{"id_token":"abc.def.ghi"}`,
			wantErr:   nil,
			wantToken: "abc.def.ghi",
		},
		{
			name:    "empty id_token",
			body:    `{"id_token":""}`,
			wantErr: ErrInvalidRequest,
		},
		{
			name:    "missing id_token field",
			body:    `{}`,
			wantErr: ErrInvalidRequest,
		},
		{
			name:    "invalid json",
			body:    `{"id_token":}`,
			wantErr: ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/login/google/firebase", strings.NewReader(tt.body))

			parsed, err := ParseRequest(req)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, parsed)
			} else {
				require.NoError(t, err)
				require.NotNil(t, parsed)
				require.Equal(t, tt.wantToken, parsed.IDToken)
			}
		})
	}
}
