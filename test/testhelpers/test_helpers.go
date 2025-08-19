package testhelpers

import (
	"context"
	"io"
	"net/http"
	"strings"

	firebaseauth "firebase.google.com/go/v4/auth"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

func NewMockSystemDeps(logger *zap.Logger) *deps.SystemDependencies {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &deps.SystemDependencies{
		MetadataURL: "https://accounts.google.com/.well-known/openid-configuration",
		HTTPClient:  http.DefaultClient,
		Logger:      logger,
	}
}

// ---- Firebase Verifier mock ----

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
	if logger == nil {
		logger = zap.NewNop()
	}
	return &deps.GoogleDependencies{
		Logger: logger,
		Verifier: &MockVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
				return &firebaseauth.Token{UID: "test-user"}, nil
			},
		},
	}
}

// ---- GitHub OAuth deps mock ----

func NewMockGitHubOAuthDeps(logger *zap.Logger) *deps.GitHubOAuthDependencies {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &deps.GitHubOAuthDependencies{
		Logger: logger,
		Config: &config.GitHubOAuthConfig{
			ClientID:    "test-client-id",
			RedirectURI: "https://example.com/callback",
			Scope:       "read:user",
			AllowSignup: "false",
		},
	}
}

// ---- HTTP client mock (implements httpclient.HTTPClient) ----

var _ httpclient.HTTPClient = (*mockHTTPClient)(nil)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)

	StatusCode int
	Body       string
	Header     http.Header
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	if m.Header == nil {
		m.Header = make(http.Header)
	}
	return &http.Response{
		StatusCode: ifZero(m.StatusCode, 200),
		Body:       io.NopCloser(strings.NewReader(m.Body)),
		Header:     m.Header,
		Request:    req,
	}, nil
}

func ifZero(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}

// ---- GitHub API deps mock ----

func NewMockGitHubAPIDeps(logger *zap.Logger) *deps.GitHubAPIDependencies {
	return NewMockGitHubAPIDepsWithClient(logger, &mockHTTPClient{
		StatusCode: 200,
		Body:       "",
	})
}

func NewMockGitHubAPIDepsWithClient(logger *zap.Logger, client httpclient.HTTPClient) *deps.GitHubAPIDependencies {
	if logger == nil {
		logger = zap.NewNop()
	}
	if client == nil {
		client = &mockHTTPClient{StatusCode: 200}
	}
	return &deps.GitHubAPIDependencies{
		APIVersion: "mock-version",
		BaseURL:    "https://api.github.com",
		HTTPClient: client,
		Logger:     logger,
		UserAgent:  "idpproxy-test",
	}
}
