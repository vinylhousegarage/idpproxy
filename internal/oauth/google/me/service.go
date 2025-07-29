package me

import (
	"net/http"
	"strings"
)

func ExtractAuthHeaderToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if strings.TrimSpace(authHeader) == "" {
		return "", ErrMissingAuthorizationHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(strings.TrimSpace(parts[0])) != "bearer" {
		return "", ErrInvalidAuthorizationHeaderFormat
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrEmptyBearerToken
	}

	return token, nil
}
