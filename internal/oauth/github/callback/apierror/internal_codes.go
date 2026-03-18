package apierror

const (
	// token exchange
	ErrorBuildRequest        ErrorCode = "build_request_failed"
	ErrorGitHubTokenRequest  ErrorCode = "github_token_request_failed"
	ErrorGitHubTokenExchange ErrorCode = "github_token_exchange_failed"

	// github user
	ErrorGitHubUserRequestBuild ErrorCode = "github_user_request_build_failed"
	ErrorGitHubUserRequest      ErrorCode = "github_user_request_failed"
	ErrorGitHubUserDecode       ErrorCode = "github_user_decode_failed"

	// internal services
	ErrorUserUpsert     ErrorCode = "user_upsert_failed"
	ErrorProxyCodeIssue ErrorCode = "proxy_code_issue_failed"
)
