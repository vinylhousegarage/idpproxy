package config

import (
	"fmt"
	"os"
)

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	ResponseType string
	Scope        string
	AccessType   string
	Prompt       string
}

func LoadGoogleConfig() (*GoogleConfig, error) {
	required := []string{
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"GOOGLE_REDIRECT_URI",
	}

	missing := []string{}
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("GoogleConfig: missing required environment variables: %q", missing)
	}

	return &GoogleConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
		ResponseType: GoogleResponseType,
		Scope:        GoogleScope,
		AccessType:   GoogleAccessType,
		Prompt:       GooglePrompt,
	}, nil
}

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
