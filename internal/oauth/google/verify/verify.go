package verify

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

type Verifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

func VerifyIDToken(ctx context.Context, client Verifier, idToken string) (*auth.Token, error) {
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verifyIDToken failed: %w", err)
	}
	return token, nil
}
