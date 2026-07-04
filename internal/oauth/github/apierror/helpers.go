package apierror

import "net/http"

// callback
func MissingGitHubCode(err error, internal ...string) *APIError {
	return New(ErrorCodeMissingGitHubCode, http.StatusBadRequest, err, internal...)
}

func MissingState(err error, internal ...string) *APIError {
	return New(ErrorCodeMissingState, http.StatusBadRequest, err, internal...)
}

func InvalidState(err error, internal ...string) *APIError {
	return New(ErrorCodeInvalidState, http.StatusBadRequest, err, internal...)
}

// token
func GitHubAccessTokenRequestError(err error, internal ...string) *APIError {
	return New(ErrorCodeGitHubAccessTokenRequest, http.StatusBadGateway, err, internal...)
}

// internal
func ServerError(err error, internal ...string) *APIError {
	return New(ErrorCodeInternal, http.StatusInternalServerError, err, internal...)
}

func ProxyCodeIssue(err error, internal ...string) *APIError {
	return New(ErrorCodeProxyCodeIssue, http.StatusInternalServerError, err, internal...)
}
