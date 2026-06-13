package apierror

const (
	// callback
	ErrorCodeMissingGitHubCode ErrorCode = "missing_github_code"
	ErrorCodeMissingState      ErrorCode = "missing_state"
	ErrorCodeInvalidState      ErrorCode = "invalid_state"

	// token
	ErrorCodeBuildAccessTokenRequest  ErrorCode = "build_access_token_request_failed"
	ErrorCodeGitHubAccessTokenRequest ErrorCode = "github_access_token_request_failed"
	ErrorCodeGitHubTokenRequest       ErrorCode = "github_token_request_failed"
	ErrorCodeGitHubTokenExchange      ErrorCode = "github_token_exchange_failed"

	// user
	ErrorCodeGitHubUserRequestBuild ErrorCode = "github_user_request_build_failed"
	ErrorCodeGitHubUserRequest      ErrorCode = "github_user_request_failed"
	ErrorCodeGitHubUserDecode       ErrorCode = "github_user_decode_failed"

	// internal
	ErrorCodeInternal       ErrorCode = "internal_error"
	ErrorCodeProxyCodeIssue ErrorCode = "proxy_code_issue_failed"
	ErrorCodeUserUpsert     ErrorCode = "user_upsert_failed"
)
