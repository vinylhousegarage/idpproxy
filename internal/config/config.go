package config

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "9000"
	}
	return port
}

func GetOpenAPIURL() string {
	url := os.Getenv("OPENAPI_URL")
	if url == "" {
		panic("OPENAPI_URL is not set")
	}
	return url
}

type FirebaseConfig struct {
	CredentialsJSON []byte
}

func LoadFirebaseConfig() (*FirebaseConfig, error) {
	b64 := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_BASE64")
	if b64 == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS_BASE64 is not set")
	}

	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode GOOGLE_APPLICATION_CREDENTIALS_BASE64: %w", err)
	}

	return &FirebaseConfig{CredentialsJSON: decoded}, nil
}

type GitHubOAuthConfig struct {
	ClientID    string
	RedirectURI string
	Scope       string
	AllowSignup string
}

func LoadGitHubOAuthConfig() (*GitHubOAuthConfig, error) {
	return loadGitHubOAuthConfigWithPrefix("GITHUB_")
}

func LoadGitHubDevOAuthConfig() (*GitHubOAuthConfig, error) {
	return loadGitHubOAuthConfigWithPrefix("GITHUB_DEV_")
}

func loadGitHubOAuthConfigWithPrefix(prefix string) (*GitHubOAuthConfig, error) {
	clientID := strings.TrimSpace(os.Getenv(prefix + "CLIENT_ID"))
	redirectURI := strings.TrimSpace(os.Getenv(prefix + "REDIRECT_URI"))

	if clientID == "" && redirectURI == "" {
		return nil, fmt.Errorf("%sCLIENT_ID and %sREDIRECT_URI are not set", prefix, prefix)
	}

	if clientID == "" {
		return nil, fmt.Errorf("%sCLIENT_ID is not set", prefix)
	}

	if redirectURI == "" {
		return nil, fmt.Errorf("%sREDIRECT_URI is not set", prefix)
	}

	if _, err := url.ParseRequestURI(redirectURI); err != nil {
		return nil, fmt.Errorf("%sREDIRECT_URI is invalid: %w", prefix, err)
	}

	return &GitHubOAuthConfig{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       GitHubScope,
		AllowSignup: GitHubAllowSignup,
	}, nil
}

type GitHubAPIConfig struct {
	APIVersion string
	BaseURL    string
	UserAgent  string
}

func LoadGitHubAPIConfig() *GitHubAPIConfig {
	return &GitHubAPIConfig{
		APIVersion: GitHubAPIVersion,
		BaseURL:    GitHubAPIBaseURL,
		UserAgent:  UserAgent(),
	}
}

type ServiceAccountConfig struct {
	ImpersonateSA string
}

func LoadServiceAccountConfig() *ServiceAccountConfig {
	return &ServiceAccountConfig{
		ImpersonateSA: strings.TrimSpace(os.Getenv("IMPERSONATE_SERVICE_ACCOUNT")),
	}
}
