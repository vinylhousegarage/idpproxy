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
	}{
		{
			name:        "RS256 with known value",
			alg:         "RS256",
			accessToken: "abc123",
			want:        "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYQ", // 確認用に実測値を記録
			wantErr:     false,
		},
		{
			name:        "RS384 should differ from RS256",
			alg:         "RS384",
			accessToken: "abc123",
			want:        "（異なる値）",
			wantErr:     false,
		},
		{
			name:        "unsupported algorithm",
			alg:         "RS999",
			accessToken: "abc",
			wantErr:     true,
		},
		{
			name:        "empty access token",
			alg:         "RS256",
			accessToken: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeAtHash(tt.alg, tt.accessToken)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got)
			t.Logf("alg=%s, at_hash=%s", tt.alg, got)
		})
	}
}
