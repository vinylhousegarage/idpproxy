package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const clientIDBytes = 32

// GenerateClientID generates a random client identifier.
func GenerateClientID() (string, error) {
	b := make([]byte, clientIDBytes)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate client ID: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
