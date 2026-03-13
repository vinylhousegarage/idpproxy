package callback

const (
	// request validation
	ErrorMissingGitHubCode = "missing_github_code"
	ErrorMissingState      = "missing_state"
	ErrorInvalidState      = "invalid_state"

	// token exchange
	ErrorBuildRequest        = "build_request_failed"
	ErrorGitHubTokenRequest  = "github_token_request_failed"
	ErrorGitHubTokenExchange = "github_token_exchange_failed"

	// github user
	ErrorGitHubUserRequestBuild = "github_user_request_build_failed"
	ErrorGitHubUserRequest      = "github_user_request_failed"
	ErrorGitHubUserDecode       = "github_user_decode_failed"

	// internal services
	ErrorUserUpsert     = "user_upsert_failed"
	ErrorProxyCodeIssue = "proxy_code_issue_failed"
)
