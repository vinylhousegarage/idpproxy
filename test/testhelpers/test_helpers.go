package testhelpers

import (
	"context"

	"go.uber.org/zap"
)

var MockLogger = zap.NewNop()

type MockVerifier struct {
	VerifyFunc func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *MockVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return m.VerifyFunc(ctx, idToken)
}
