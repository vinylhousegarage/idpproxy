package apierror

const (
	// request validation
	ErrorMissingGitHubCode ErrorCode = "missing_github_code"
	ErrorMissingState      ErrorCode = "missing_state"
	ErrorInvalidState      ErrorCode = "invalid_state"
)
