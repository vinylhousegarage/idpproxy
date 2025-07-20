package callback

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/apperror"
)

var (
	ErrFailedToCreateRequest        = apperror.New(http.StatusInternalServerError, "failed to create request")        // 500 Internal Server Error
	ErrFailedToDecodeTokenResponse  = apperror.New(http.StatusInternalServerError, "failed to decode token response") // 500 Internal Server Error
	ErrFailedToDecodeMetadata       = apperror.New(http.StatusBadGateway, "failed to decode metadata")                // 502 Bad Gateway
	ErrFailedToEncodeTokenResponse  = apperror.New(http.StatusInternalServerError, "failed to encode token response") // 500 Internal Server Error
	ErrFailedToFetchMetadata        = apperror.New(http.StatusBadGateway, "failed to fetch metadata")                 // 502 Bad Gateway
	ErrInvalidState                 = apperror.New(http.StatusBadRequest, "invalid state")                            // 400 Bad Request
	ErrMissingCode                  = apperror.New(http.StatusBadRequest, "missing code")                             // 400 Bad Request
	ErrMissingState                 = apperror.New(http.StatusBadRequest, "missing state")                            // 400 Bad Request
	ErrMissingStateCookie           = apperror.New(http.StatusBadRequest, "missing oauth_state cookie")               // 400 Bad Request
	ErrMissingTokenEndpoint         = apperror.New(http.StatusBadGateway, "missing token_endpoint")                   // 502 Bad Gateway
	ErrUnexpectedCallbackStatusCode = apperror.New(http.StatusBadGateway, "unexpected response from callback")        // 502 Bad Gateway
	ErrUnexpectedMetadataStatusCode = apperror.New(http.StatusBadGateway, "unexpected response from metadata")        // 502 Bad Gateway
)
