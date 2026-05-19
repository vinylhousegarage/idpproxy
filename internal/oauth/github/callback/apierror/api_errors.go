package apierror

import "net/http"

func MissingGitHubCode(err error, internal ...string) *APIError {
	return New(ErrorMissingGitHubCode, http.StatusBadRequest, err, internal...)
}

func MissingState(err error, internal ...string) *APIError {
	return New(ErrorMissingState, http.StatusBadRequest, err, internal...)
}

func InvalidState(err error, internal ...string) *APIError {
	return New(ErrorInvalidState, http.StatusBadRequest, err, internal...)
}

func ProxyCodeIssue(err error, internal ...string) *APIError {
	return New(ErrorProxyCodeIssue, http.StatusInternalServerError, err, internal...)
}

func Internal(err error, internal ...string) *APIError {
	return New(ErrorInternal, http.StatusInternalServerError, err, internal...)
}
