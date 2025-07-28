package me

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/apperror"
)

var (
	ErrInvalidAuthorizationHeaderFormat = apperror.New(http.StatusBadRequest, "invalid authorization header format") // 400 Bad Request
	ErrMissingAuthorizationHeader       = apperror.New(http.StatusUnauthorized, "missing authorization header")      // 401 Unauthorized
)
