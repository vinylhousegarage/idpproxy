package apierror

import "net/http"

const (
	ErrorCodeInternal       ErrorCode = "internal_error"
	ErrorCodeProxyCodeIssue ErrorCode = "proxy_code_issue_failed"
	ErrorCodeUserUpsert     ErrorCode = "user_upsert_failed"
)

func Internal(err error, internal ...string) *APIError {
	return New(ErrorCodeInternal, http.StatusInternalServerError, err, internal...)
}

func ProxyCodeIssue(err error, internal ...string) *APIError {
	return New(ErrorCodeProxyCodeIssue, http.StatusInternalServerError, err, internal...)
}
