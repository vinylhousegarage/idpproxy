package config

import (
	"encoding/base64"
	"fmt"
	"os"
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

type GitHubConfig struct {
	ClientID    string
	RedirectURI string
	Scope       string
	AllowSignup string
}

func LoadGitHubConfig() (*GitHubConfig, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	if clientID == "" || redirectURI == "" {
		return nil, fmt.Errorf("missing GitHub client ID or redirect URI")
	}

	return &GitHubConfig{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       "read:user",
		AllowSignup: "true",
	}, nil
}
func LoadGitHubDevConfig() (*GitHubConfig, error) {
	clientID := os.Getenv("GITHUB_DEV_CLIENT_ID")
	redirectURI := os.Getenv("GITHUB_DEV_REDIRECT_URI")
	if clientID == "" || redirectURI == "" {
		return nil, fmt.Errorf("missing GitHub client ID or redirect URI")
	}

	return &GitHubConfig{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       "read:user",
		AllowSignup: "true",
	}, nil
}
