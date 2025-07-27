package testhelpers

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func NewMockSystemDeps(logger *zap.Logger) *deps.SystemDependencies {
	return &deps.SystemDependencies{
		MetadataURL: "https://accounts.google.com/.well-known/openid-configuration",
		HTTPClient:  http.DefaultClient,
		Logger:      logger,
	}
}

type MockVerifier struct {
	VerifyFunc func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (m *MockVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return m.VerifyFunc(ctx, idToken)
}

func NewMockGoogleDeps(logger *zap.Logger) *deps.GoogleDependencies {
	return &deps.GoogleDependencies{
		Logger: logger,
		Verifier: &MockVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{UID: "test-user"}, nil
			},
		},
	}
}
