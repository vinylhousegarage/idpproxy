package callback

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback/apierror"
	"go.uber.org/zap"
)

var logger = zap.NewNop()

func writeJSON(w http.ResponseWriter, status int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error("writeJSON marshal error", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		logger.Error("writeJSON write error", zap.Error(err))
	}
}

func WriteError(c *gin.Context, err error) {
	var apiErr *apierror.APIError

	if ok := asAPIError(err, &apiErr); ok {
		writeJSON(c.Writer, apiErr.HTTPStatus, apierror.ErrorResponse{
			Error: apiErr.Code,
		})

		return
	}

	internal := apierror.Internal(err)

	writeJSON(c.Writer, internal.HTTPStatus, apierror.ErrorResponse{
		Error: internal.Code,
	})
}

func asAPIError(err error, target **apierror.APIError) bool {
	return errorsAs(err, target)
}

var errorsAs = func(err error, target any) bool {
	return errors.As(err, target)
}
