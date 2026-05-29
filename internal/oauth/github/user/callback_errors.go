package user

import "github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"

const (
	ErrorCodeGitHubUserRequestBuild apierror.ErrorCode = "github_user_request_build_failed"
	ErrorCodeGitHubUserRequest      apierror.ErrorCode = "github_user_request_failed"
	ErrorCodeGitHubUserDecode       apierror.ErrorCode = "github_user_decode_failed"
)
