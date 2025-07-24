package verify

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/require"
)

type mockAuthClient struct {
	mockVerify func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *mockAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return m.mockVerify(ctx, idToken)
}

func TestVerifyIDToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("valid token returns token object", func(t *testing.T) {
		t.Parallel()

		mock := &mockAuthClient{
			mockVerify: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{
					UID:    "user123",
					Claims: map[string]interface{}{"email": "test@example.com"},
				}, nil
			},
		}

		token, err := VerifyIDToken(ctx, mock, "valid_token")
		require.NoError(t, err)
		require.Equal(t, "user123", token.UID)
		require.Equal(t, "test@example.com", token.Claims["email"])
	})

	t.Run("invalid token returns error", func(t *testing.T) {
		t.Parallel()

		mock := &mockAuthClient{
			mockVerify: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return nil, errors.New("invalid token")
			},
		}

		token, err := VerifyIDToken(ctx, mock, "invalid_token")
		require.Error(t, err)
		require.Nil(t, token)
		require.Contains(t, err.Error(), "invalid token")
	})
}
