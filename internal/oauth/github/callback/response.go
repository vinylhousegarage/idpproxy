package callback

import "github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"

type ErrorResponse struct {
	Error apierror.ErrorCode `json:"error"`
}
