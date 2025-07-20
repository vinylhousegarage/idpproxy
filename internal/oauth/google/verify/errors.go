package verify

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/apperror"
)

var (
	ErrInvalidAudience   = apperror.New(http.StatusUnauthorized, "invalid audience")           // 401 Unauthorized
	ErrInvalidIssuer     = apperror.New(http.StatusUnauthorized, "unexpected issuer")          // 401 Unauthorized
	ErrInvalidSigningAlg = apperror.New(http.StatusBadRequest, "unexpected signing method")    // 400 Bad Request
	ErrJWTParseFailed    = apperror.New(http.StatusInternalServerError, "failed to parse JWT") // 500 Internal Server Error
	ErrMissingAudience   = apperror.New(http.StatusUnauthorized, "audience claim is missing")  // 401 Unauthorized
	ErrMissingSubject    = apperror.New(http.StatusUnauthorized, "missing subject (sub)")      // 401 Unauthorized
	ErrTokenExpired      = apperror.New(http.StatusUnauthorized, "token is expired")           // 401 Unauthorized
)
