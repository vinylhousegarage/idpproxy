package signer

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func makeHS256Token(t *testing.T, key []byte, kid, typ string, claims jwt.MapClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if kid != "" {
		token.Header["kid"] = kid
	}
	if typ == "" {
		delete(token.Header, "typ")
	} else {
		token.Header["typ"] = typ
	}

	signed, err := token.SignedString(key)
	require.NoError(t, err)
	require.NotEmpty(t, signed)

	return signed
}

func TestHMACSigner_Verify(t *testing.T) {
	t.Parallel()

	secret := []byte("secret")
	kid := "kid"
	base := time.Unix(1_700_000_000, 0)

	tests := []struct {
		name      string
		claims    jwt.MapClaims
		typ       string
		kidHeader string
		opt       *VerifyOptions
		now       time.Time
		key       []byte
		wantErr   error
	}{
		{
			name: "basic",
			claims: jwt.MapClaims{
				"sub": "u1",
				"iat": base.Unix(),
				"nbf": base.Unix(),
				"exp": base.Add(10 * time.Minute).Unix(),
			},
			typ:       "JWT",
			kidHeader: kid,
			opt:       &VerifyOptions{},
			now:       base,
			key:       secret,
			wantErr:   nil,
		},
		{
			name: "require typ ok",
			claims: jwt.MapClaims{
				"exp": base.Add(5 * time.Minute).Unix(),
			},
			typ:       "JWT",
			kidHeader: kid,
			opt:       &VerifyOptions{RequireTyp: true},
			now:       base,
			key:       secret,
			wantErr:   nil,
		},
		{
			name: "require typ missing",
			claims: jwt.MapClaims{
				"exp": base.Add(5 * time.Minute).Unix(),
			},
			typ:       "",
			kidHeader: kid,
			opt:       &VerifyOptions{RequireTyp: true},
			now:       base,
			key:       secret,
			wantErr:   ErrInvalidTyp,
		},
		{
			name: "expect kid mismatch",
			claims: jwt.MapClaims{
				"exp": base.Add(5 * time.Minute).Unix(),
			},
			typ:       "JWT",
			kidHeader: "other-kid",
			opt:       &VerifyOptions{ExpectKID: "expected"},
			now:       base,
			key:       secret,
			wantErr:   ErrUnexpectedKID,
		},
		{
			name: "expired no leeway",
			claims: jwt.MapClaims{
				"exp": base.Add(-10 * time.Second).Unix(),
			},
			typ:       "JWT",
			kidHeader: kid,
			opt:       &VerifyOptions{},
			now:       base,
			key:       secret,
			wantErr:   jwt.ErrTokenExpired,
		},
		{
			name: "expired with leeway ok",
			claims: jwt.MapClaims{
				"exp": base.Add(10 * time.Second).Unix(),
			},
			typ:       "JWT",
			kidHeader: kid,
			opt:       &VerifyOptions{Leeway: 30 * time.Second},
			now:       base.Add(30 * time.Second),
			key:       secret,
			wantErr:   nil,
		},
		{
			name: "not before with leeway ok",
			claims: jwt.MapClaims{
				"nbf": base.Add(10 * time.Second).Unix(),
				"exp": base.Add(5 * time.Minute).Unix(),
			},
			typ:       "JWT",
			kidHeader: kid,
			opt:       &VerifyOptions{Leeway: 10 * time.Second},
			now:       base.Add(5 * time.Second),
			key:       secret,
			wantErr:   nil,
		},
		{
			name:    "empty token",
			claims:  nil,
			opt:     nil,
			now:     base,
			key:     secret,
			wantErr: ErrEmptyToken,
		},
		{
			name:    "empty key",
			claims:  nil,
			opt:     nil,
			now:     base,
			key:     nil,
			wantErr: ErrEmptyKey,
		},
		{
			name:      "expect kid set but header missing",
			claims:    jwt.MapClaims{"exp": base.Add(5 * time.Minute).Unix()},
			typ:       "JWT",
			kidHeader: "",
			opt:       &VerifyOptions{ExpectKID: "expected-kid"},
			now:       base,
			key:       secret,
			wantErr:   ErrUnexpectedKID,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewHMACSigner(tt.key, kid)
			s.now = func() time.Time { return tt.now }

			var token string
			if tt.claims != nil {
				token = makeHS256Token(t, tt.key, tt.kidHeader, tt.typ, tt.claims)
			}

			got, err := s.Verify(context.Background(), token, tt.opt)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
