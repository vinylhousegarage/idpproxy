package callback

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
)

const (
	ErrorCodeMissingGitHubCode apierror.ErrorCode = "missing_github_code"
	ErrorCodeMissingState      apierror.ErrorCode = "missing_state"
	ErrorCodeInvalidState      apierror.ErrorCode = "invalid_state"
)

func MissingGitHubCode(err error, internal ...string) *apierror.APIError {
	return apierror.New(ErrorCodeMissingGitHubCode, http.StatusBadRequest, err, internal...)
}

func MissingState(err error, internal ...string) *apierror.APIError {
	return apierror.New(ErrorCodeMissingState, http.StatusBadRequest, err, internal...)
}

func InvalidState(err error, internal ...string) *apierror.APIError {
	return apierror.New(ErrorCodeInvalidState, http.StatusBadRequest, err, internal...)
}
