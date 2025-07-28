package me

import (
	"net/http"
	"strings"
)

func ExtractAuthHeaderToken(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", ErrMissingAuthorizationHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidAuthorizationHeaderFormat
	}

	return parts[1], nil
}
