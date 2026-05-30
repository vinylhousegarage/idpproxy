package apierror

import "net/http"

func ProxyCodeIssue(err error, internal ...string) *APIError {
	return New(ErrorCodeProxyCodeIssue, http.StatusInternalServerError, err, internal...)
}

func Internal(err error, internal ...string) *APIError {
	return New(ErrorCodeInternal, http.StatusInternalServerError, err, internal...)
}
