package apierror

import (
	"errors"
)

func MissingState(err error) *APIError {
	return New(ErrorMissingState, 400, err)
}

func MissingGitHubCode(err error) *APIError {
	return New(ErrorMissingGitHubCode, 400, err)
}

func InvalidState(err error) *APIError {
	return New(ErrorInvalidState, 400, err)
}

func Internal(err error) *APIError {
	return New(ErrorInternal, 500, err)
}

func FromInternal(err error) *APIError {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return InvalidState(err)

	case errors.Is(err, ErrPersist):
		return New(ErrorProxyCodeIssue, 500, err)

	default:
		return Internal(err)
	}
}
