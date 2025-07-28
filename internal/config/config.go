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
	if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials from file: %w", err)
		}
		return &FirebaseConfig{CredentialsJSON: data}, nil
	}

	if b64 := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_BASE64"); b64 != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode GOOGLE_APPLICATION_CREDENTIALS_BASE64: %w", err)
		}
		return &FirebaseConfig{CredentialsJSON: decoded}, nil
	}

	return nil, fmt.Errorf("no Firebase credentials found in environment")
}
