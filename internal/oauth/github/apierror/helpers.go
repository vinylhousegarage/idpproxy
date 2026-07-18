package apierror

import "net/http"

// callback
func MissingGitHubCode(err error, internals ...APIInternal) *APIError {
	return New(ErrorCodeMissingGitHubCode, http.StatusBadRequest, err, internals...)
}

func MissingState(err error, internals ...APIInternal) *APIError {
	return New(ErrorCodeMissingState, http.StatusBadRequest, err, internals...)
}

func InvalidState(err error, internals ...APIInternal) *APIError {
	return New(ErrorCodeInvalidState, http.StatusBadRequest, err, internals...)
}

// token
func GitHubAccessTokenRequestError(err error, internals ...APIInternal) *APIError {
	return New(ErrorCodeGitHubAccessTokenRequest, http.StatusBadGateway, err, internals...)
}

// internal
func InternalServerError(err error, internals ...APIInternal) *APIError {
	return New(ErrorCodeInternalServerError, http.StatusInternalServerError, err, internals...)
}

func ProxyCodeIssue(err error, internals ...APIInternal) *APIError {
	return New(ErrorCodeProxyCodeIssue, http.StatusInternalServerError, err, internals...)
}
