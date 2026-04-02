package apierror

import "net/http"

func MissingGitHubCode(err error) *APIError {
	return New(ErrorMissingGitHubCode, http.StatusBadRequest, err)
}

func MissingState(err error) *APIError {
	return New(ErrorMissingState, http.StatusBadRequest, err)
}

func InvalidState(err error) *APIError {
	return New(ErrorInvalidState, http.StatusUnauthorized, err)
}

func ProxyCodeIssue(err error) *APIError {
	return New(ErrorProxyCodeIssue, http.StatusInternalServerError, err)
}

func Internal(err error) *APIError {
	return New(ErrorInternal, http.StatusInternalServerError, err)
}
