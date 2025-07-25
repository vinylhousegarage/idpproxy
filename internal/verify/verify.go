package verify

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type Verifier struct {
	Auth *auth.Client
}

func NewVerifier(authClient *auth.Client) *Verifier {
	return &Verifier{Auth: authClient}
}

func (v *Verifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := v.Auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired id_token: %w", err)
	}

	return token, nil
}
