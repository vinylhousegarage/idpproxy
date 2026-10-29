package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

type IDTokenConfig struct {
	Issuer    string
	HMACKey   []byte
	HMACKeyID string
	TTL       time.Duration
}

func LoadIDTokenConfig() (*IDTokenConfig, error) {
	issuer := strings.TrimSpace(os.Getenv("ID_TOKEN_ISSUER"))
	if issuer == "" {
		return nil, fmt.Errorf("ID_TOKEN_ISSUER is not set")
	}

	hmacKeyBase64 := strings.TrimSpace(
		os.Getenv("ID_TOKEN_HMAC_KEY_BASE64"),
	)
	if hmacKeyBase64 == "" {
		return nil, fmt.Errorf("ID_TOKEN_HMAC_KEY_BASE64 is not set")
	}

	hmacKey, err := base64.StdEncoding.DecodeString(hmacKeyBase64)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode ID_TOKEN_HMAC_KEY_BASE64: %w",
			err,
		)
	}

	hmacKeyID := strings.TrimSpace(os.Getenv("ID_TOKEN_HMAC_KEY_ID"))
	if hmacKeyID == "" {
		return nil, fmt.Errorf("ID_TOKEN_HMAC_KEY_ID is not set")
	}

	ttlValue := strings.TrimSpace(os.Getenv("ID_TOKEN_TTL"))
	if ttlValue == "" {
		return nil, fmt.Errorf("ID_TOKEN_TTL is not set")
	}

	ttl, err := time.ParseDuration(ttlValue)
	if err != nil {
		return nil, fmt.Errorf("ID_TOKEN_TTL is invalid: %w", err)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("ID_TOKEN_TTL must be greater than zero")
	}

	return &IDTokenConfig{
		Issuer:    issuer,
		HMACKey:   hmacKey,
		HMACKeyID: hmacKeyID,
		TTL:       ttl,
	}, nil
}
