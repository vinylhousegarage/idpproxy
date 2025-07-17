package httperror

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/apperror"

	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error" example:"invalid token"`
}

func WriteJSONError(w http.ResponseWriter, status int, msg string, logger *zap.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: msg}); err != nil {
		logger.Error("failed to write error response", zap.Error(err))
	}
}

func WriteErrorResponse(w http.ResponseWriter, err error, logger *zap.Logger) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		WriteJSONError(w, appErr.StatusCode(), appErr.Error(), logger)
	} else {
		logger.Error("unhandled internal error", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, "internal server error", logger)
	}
}
