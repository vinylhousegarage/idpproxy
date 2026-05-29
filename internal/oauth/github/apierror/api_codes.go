package apierror

const (
	// validation
	ErrorCodeMissingGitHubCode ErrorCode = "missing_github_code"
	ErrorCodeMissingState      ErrorCode = "missing_state"
	ErrorCodeInvalidState      ErrorCode = "invalid_state"

	// internal services
	ErrorCodeUserUpsert     ErrorCode = "user_upsert_failed"
	ErrorCodeProxyCodeIssue ErrorCode = "proxy_code_issue_failed"

	// fallback
	ErrorCodeInternal ErrorCode = "internal_error"
)
