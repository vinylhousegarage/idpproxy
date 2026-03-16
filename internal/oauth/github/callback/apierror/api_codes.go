package apierror

import "net/http"

var (
	ErrMissingGitHubCode = New(ErrorMissingGitHubCode, http.StatusBadRequest, nil)
	ErrMissingState      = New(ErrorMissingState, http.StatusBadRequest, nil)
	ErrInvalidState      = New(ErrorInvalidState, http.StatusUnauthorized, nil)
)
