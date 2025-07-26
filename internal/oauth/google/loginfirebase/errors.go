package loginfirebase

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/apperror"
)

var (
	ErrInvalidIDToken = apperror.New(http.StatusUnauthorized, "invalid id_token") // 401 Unauthorized
	ErrInvalidRequest = apperror.New(http.StatusBadRequest, "invalid request")    // 400 Bad Request
)
