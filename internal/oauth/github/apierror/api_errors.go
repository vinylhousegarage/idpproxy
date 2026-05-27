package apierror

import "net/http"

func MissingGitHubCode(err error, internal ...string) *APIError {
	return New(ErrorCodeMissingGitHubCode, http.StatusBadRequest, err, internal...)
}

func MissingState(err error, internal ...string) *APIError {
	return New(ErrorCodeMissingState, http.StatusBadRequest, err, internal...)
}

func InvalidState(err error, internal ...string) *APIError {
	return New(ErrorCodeInvalidState, http.StatusBadRequest, err, internal...)
}

func ProxyCodeIssue(err error, internal ...string) *APIError {
	return New(ErrorCodeProxyCodeIssue, http.StatusInternalServerError, err, internal...)
}

func Internal(err error, internal ...string) *APIError {
	return New(ErrorCodeInternal, http.StatusInternalServerError, err, internal...)
}
