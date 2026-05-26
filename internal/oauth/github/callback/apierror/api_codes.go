package apierror

const (
	// validation
	ErrorCodeMissingGitHubCode ErrorCode = "missing_github_code"
	ErrorCodeMissingState      ErrorCode = "missing_state"
	ErrorCodeInvalidState      ErrorCode = "invalid_state"

	// github token
	ErrorCodeBuildRequest        ErrorCode = "build_request_failed"
	ErrorCodeGitHubTokenRequest  ErrorCode = "github_token_request_failed"
	ErrorCodeGitHubTokenExchange ErrorCode = "github_token_exchange_failed"

	// github user
	ErrorCodeGitHubUserRequestBuild ErrorCode = "github_user_request_build_failed"
	ErrorCodeGitHubUserRequest      ErrorCode = "github_user_request_failed"
	ErrorCodeGitHubUserDecode       ErrorCode = "github_user_decode_failed"

	// internal services
	ErrorCodeUserUpsert     ErrorCode = "user_upsert_failed"
	ErrorCodeProxyCodeIssue ErrorCode = "proxy_code_issue_failed"

	// fallback
	ErrorCodeInternal ErrorCode = "internal_error"
)
