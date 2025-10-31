package idtoken

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeAtHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		alg         string
		accessToken string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:        "RS256 with known value",
			alg:         "RS256",
			accessToken: "abc123",
			want:        "bKE9UspwyIPg8LsQHkJaiQ",
			wantErr:     false,
		},
		{
			name:        "RS384 with known value",
			alg:         "RS384",
			accessToken: "abc123",
			want:        "ox15iRkZytJPMmRHnXaIT1gb7jLoZ3g3",
			wantErr:     false,
		},
		{
			name:        "RS512 with known value",
			alg:         "RS512",
			accessToken: "abc123",
			want:        "xwtd2ev7b1HQnUEytxcMnSB1CnhS8AaA9lZY8DEOgQA",
			wantErr:     false,
		},
		{
			name:        "lowercase alg is accepted",
			alg:         "rs256",
			accessToken: "abc123",
			want:        "bKE9UspwyIPg8LsQHkJaiQ",
			wantErr:     false,
		},
		{
			name:        "unsupported algorithm",
			alg:         "RS999",
			accessToken: "abc",
			wantErr:     true,
			errContains: "unsupported alg",
		},
		{
			name:        "empty access token",
			alg:         "RS256",
			accessToken: "",
			wantErr:     true,
			errContains: "empty access token",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := computeAtHash(tc.alg, tc.accessToken)
			if tc.wantErr {
				require.Error(t, err)
				if tc.errContains != "" {
					require.ErrorContains(t, err, tc.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got, "alg=%s token=%q", tc.alg, tc.accessToken)
		})
	}
}
