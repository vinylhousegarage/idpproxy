package idtoken

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIDTokenClaims_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		claims  IDTokenClaims
		wantErr error
	}{
		{
			name: "ok: minimal required fields",
			claims: IDTokenClaims{
				Iss: "https://idpproxy.com",
				Sub: "user-123",
				Aud: "client-abc",
				Iat: 1_800_000_000,
				Exp: 1_800_000_600,
			},
			wantErr: nil,
		},
		{
			name: "ng: empty iss",
			claims: IDTokenClaims{
				Sub: "user-123", Aud: "client-abc",
				Iat: 1_800_000_000, Exp: 1_800_000_600,
			},
			wantErr: ErrInvalidIssuer,
		},
		{
			name: "ng: empty sub",
			claims: IDTokenClaims{
				Iss: "https://idpproxy.com", Aud: "client-abc",
				Iat: 1_800_000_000, Exp: 1_800_000_600,
			},
			wantErr: ErrInvalidSubject,
		},
		{
			name: "ng: empty aud",
			claims: IDTokenClaims{
				Iss: "https://idpproxy.com", Sub: "user-123",
				Iat: 1_800_000_000, Exp: 1_800_000_600,
			},
			wantErr: ErrInvalidAudience,
		},
		{
			name: "ng: iat <= 0",
			claims: IDTokenClaims{
				Iss: "https://idpproxy.com", Sub: "user-123", Aud: "client-abc",
				Iat: 0, Exp: 1_800_000_600,
			},
			wantErr: ErrInvalidIat,
		},
		{
			name: "ng: exp == iat",
			claims: IDTokenClaims{
				Iss: "https://idpproxy.com", Sub: "user-123", Aud: "client-abc",
				Iat: 1_800_000_000, Exp: 1_800_000_000,
			},
			wantErr: ErrInvalidExp,
		},
		{
			name: "ng: exp < iat",
			claims: IDTokenClaims{
				Iss: "https://idpproxy.com", Sub: "user-123", Aud: "client-abc",
				Iat: 1_800_000_000, Exp: 1_799_999_999,
			},
			wantErr: ErrInvalidExp,
		},
		{
			name: "ok: optional fields present do not affect validation",
			claims: IDTokenClaims{
				Iss:      "https://idpproxy.com",
				Sub:      "user-123",
				Aud:      "client-abc",
				Iat:      1_800_000_000,
				Exp:      1_800_000_600,
				Nonce:    "n-xyz",
				AuthTime: 1_799_999_500,
				AMR:      []string{"pwd", "mfa"},
				AtHash:   "abcd",
				Azp:      "client-abc",
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.claims.Validate()
			if tc.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}
