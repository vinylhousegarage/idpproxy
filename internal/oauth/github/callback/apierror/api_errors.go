package apierror

const (
	// validation errors
	ErrorMissingGitHubCode ErrorCode = "missing_github_code"
	ErrorMissingState      ErrorCode = "missing_state"
	ErrorInvalidState      ErrorCode = "invalid_state"

	// internal errors
	ErrorInternal ErrorCode = "internal_error"
)
