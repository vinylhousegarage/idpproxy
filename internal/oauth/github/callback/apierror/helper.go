package apierror

import "errors"

func FromInternal(err error) *APIError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrInvalidInput):
		return InvalidState(err)

	case errors.Is(err, ErrPersist):
		return ProxyCodeIssue(err)

	default:
		return Internal(err)
	}
}
