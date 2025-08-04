package testhelpers

import (
	"context"
	"net/http"

	firebaseauth "firebase.google.com/go/v4/auth"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
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
	VerifyFunc func(ctx context.Context, idToken string) (*firebaseauth.Token, error)
}

func (m *MockVerifier) VerifyIDToken(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
	if m.VerifyFunc != nil {
		return m.VerifyFunc(ctx, idToken)
	}
	return &firebaseauth.Token{UID: "default-mock"}, nil
}

func NewMockGoogleDeps(logger *zap.Logger) *deps.GoogleDependencies {
	return &deps.GoogleDependencies{
		Logger: logger,
		Verifier: &MockVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
				return &firebaseauth.Token{UID: "test-user"}, nil
			},
		},
	}
}

func NewMockGitHubDeps(logger *zap.Logger) *deps.GitHubDependencies {
    return &deps.GitHubDependencies{
        Logger: logger,
        Config: &config.GitHubConfig{
            ClientID:    "dummy-client-id",
            RedirectURI: "https://example.com/callback",
            Scope:       "read:user",
            AllowSignup: "false",
        },
    }
}
